package handler

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	configure "github.com/livesense-inc/fanlin/lib/conf"
	"github.com/livesense-inc/fanlin/lib/content"
	_ "github.com/livesense-inc/fanlin/lib/content/local"
	_ "github.com/livesense-inc/fanlin/lib/content/s3"
	_ "github.com/livesense-inc/fanlin/lib/content/web"
	imgproxyerr "github.com/livesense-inc/fanlin/lib/error"
	imageprocessor "github.com/livesense-inc/fanlin/lib/image"
	"github.com/livesense-inc/fanlin/lib/query"
	_ "github.com/livesense-inc/fanlin/plugin"
	servertiming "github.com/mitchellh/go-server-timing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

var devNull, _ = os.Open("/dev/null")

func create404Page(w http.ResponseWriter, r *http.Request, conf *configure.Conf) {
	q := query.NewQueryFromGet(r)
	width, height := clampBounds(conf, q)
	w.WriteHeader(http.StatusNotFound)
	var b bytes.Buffer
	if err := imageprocessor.Set404Image(
		&b,
		content.GetNoContentImage(),
		width,
		height,
		*q.FillColor(),
	); err != nil {
		writeDebugLog(err, conf.DebugLogPath())
		log.Println(err)
		fmt.Fprintf(w, "%s", "404 Not found.")
	} else {
		_, _ = io.Copy(w, &b)
	}
}

func fallback(
	w http.ResponseWriter,
	r *http.Request,
	conf *configure.Conf,
	loggers map[string]*logrus.Logger,
	err error,
) {
	create404Page(w, r, conf)
	if err == nil {
		return
	}
	if loggers != nil {
		errLogger := func() *logrus.Entry {
			logger := loggers["err"]
			return logger.WithFields(logrus.Fields{
				"UA":        r.UserAgent(),
				"access_ip": r.RemoteAddr,
				"url":       r.URL.String(),
				"type":      r.Method,
				"version":   r.Proto,
			})
		}()
		if e, ok := err.(*imgproxyerr.Err); ok {
			switch e.Type {
			case imgproxyerr.WARNING:
				os.Stderr = devNull
				errLogger.Warn(err)
			case imgproxyerr.ERROR:
				writeDebugLog(err, conf.DebugLogPath())
				errLogger.Error(err)
			default:
				writeDebugLog(err, conf.DebugLogPath())
				errLogger.Error(err)
			}
		} else {
			writeDebugLog(err, conf.DebugLogPath())
			errLogger.Error(err)
		}
	} else {
		writeDebugLog(err, conf.DebugLogPath())
		log.Println(err)
	}
	fmt.Fprintf(w, "%s", "")
}

func writeDebugLog(err interface{}, debugFile string) {
	stackWriter, _ := os.OpenFile(debugFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	t := time.Now()
	stackWriter.Write([]byte("\n"))
	stackWriter.Write([]byte("==========================================\n"))
	stackWriter.Write([]byte(t.String() + "\n"))
	stackWriter.Write([]byte(fmt.Sprint(err, "\n")))
	stackWriter.Write([]byte("==========================================\n\n"))
	os.Stderr = stackWriter
}

func MainHandler(
	w http.ResponseWriter,
	r *http.Request,
	conf *configure.Conf,
	loggers map[string]*logrus.Logger,
) {
	accessLogger := loggers["access"]
	accessLogger.WithFields(logrus.Fields{
		"UA":        r.UserAgent(),
		"access_ip": r.RemoteAddr,
		"url":       r.URL.String(),
	}).Info()

	timing := servertiming.FromContext(r.Context())

	m := timing.NewMetric("f_load").Start()
	buf, err := getImage(r.URL.Path, conf)
	if err != nil {
		fallback(
			w, r, conf, loggers,
			imgproxyerr.New(imgproxyerr.WARNING, fmt.Errorf("failed to get image data: %w", err)),
		)
		return
	}
	if buf == nil {
		create404Page(w, r, conf)
		return
	}
	m.Stop()

	q := query.NewQueryFromGet(r)

	m = timing.NewMetric("f_process").Start()
	img, err := processImage(buf, conf, q)
	if err != nil {
		fallback(
			w, r, conf, loggers,
			imgproxyerr.New(imgproxyerr.ERROR, fmt.Errorf("failed to decode image data: %w", err)),
		)
		return
	}
	m.Stop()

	m = timing.NewMetric("f_encode").Start()
	var b bytes.Buffer
	if err := encodeImage(&b, img, q); err != nil {
		fallback(
			w, r, conf, loggers,
			imgproxyerr.New(imgproxyerr.ERROR, fmt.Errorf("failed to encode image data: %w", err)),
		)
		return
	}
	m.Stop()

	if q.UseAVIF() {
		w.Header().Set("Content-Type", "image/avif")
	}
	w.WriteHeader(http.StatusOK)
	if _, err := io.Copy(w, &b); err != nil {
		loggers["err"].Print("failed to write data to response:", err)
	}
}

func getImage(reqPath string, conf *configure.Conf) (io.Reader, error) {
	ctt := content.GetContent(reqPath, conf)
	if ctt == nil {
		return nil, nil
	}
	return content.GetImageBinary(ctt)
}

func processImage(buf io.Reader, conf *configure.Conf, q *query.Query) (*imageprocessor.Image, error) {
	img, err := imageprocessor.DecodeImage(buf)
	if err != nil {
		return nil, err
	}
	if conf.UseMLCMYKConverter() {
		if err := img.ConvertColor(conf.MLCMYKConverterNetworkFilePath()); err != nil {
			return nil, err
		}
	} else if conf.UseICCProfileCMYKConverter() {
		img.ConvertColorWithICCProfile()
	}
	img.ApplyOrientation()
	w, h := clampBounds(conf, q)
	if q.Crop() {
		img.Crop(w, h)
	} else {
		img.ResizeAndFill(w, h, *q.FillColor())
	}
	if q.Grayscale() {
		img.Grayscale()
	} else if q.Inverse() {
		img.Invert()
	}
	img.Process()
	return img, nil
}

func clampBounds(conf *configure.Conf, q *query.Query) (w uint, h uint) {
	mW, mX := conf.MaxSize()
	b := q.Bounds()
	w = min(b.W, mW)
	h = min(b.H, mX)
	return
}

func encodeImage(
	w io.Writer,
	img *imageprocessor.Image,
	q *query.Query,
) (err error) {
	switch img.GetFormat() {
	case "jpeg":
		if q.UseWebP() {
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), false)
		} else if q.UseAVIF() {
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodeJpeg(w, img.GetImg(), q.Quality())
		}
	case "png":
		if q.UseWebP() {
			useLossless := (q.Quality() == 100)
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), useLossless)
		} else if q.UseAVIF() {
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodePNG(w, img.GetImg(), q.Quality())
		}
	case "gif":
		if q.UseWebP() {
			useLossless := (q.Quality() == 100)
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), useLossless)
		} else if q.UseAVIF() {
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodeGIF(w, img.GetImg(), q.Quality())
		}
	case "webp":
		useLossless := (q.Quality() == 100)
		err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), useLossless)
	case "avif":
		err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
	default:
		if q.UseWebP() {
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), false)
		} else if q.UseAVIF() {
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodeJpeg(w, img.GetImg(), q.Quality())
		}
	}
	return
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", "")
}

func MakeMetricsHandler(conf *configure.Conf, logger *log.Logger) http.Handler {
	return promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				DisableCompression: true,
				ErrorLog:           logger,
				Timeout:            conf.BackendRequestTimeout(),
			},
		),
	)
}

func Prepare(conf *configure.Conf) error {
	content.SetUpProviders(conf)
	if err := content.SetupNoContentImage(conf); err != nil {
		return err
	}
	if err := imageprocessor.SetUpColorConverter(); err != nil {
		return err
	}
	return nil
}
