package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
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

	maxW, maxH := conf.MaxSize()
	w.WriteHeader(404)
	if err := imageprocessor.Set404Image(w, conf.NotFoundImagePath(), q.Bounds().W, q.Bounds().H, *q.FillColor(), maxW, maxH); err != nil {
		writeDebugLog(err, conf.DebugLogPath())
		log.Println(err)
		fmt.Fprintf(w, "%s", "404 Not found.")
	}

	q = nil
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

func MainHandler(w http.ResponseWriter, r *http.Request, conf *configure.Conf, loggers map[string]*logrus.Logger) {
	timing := servertiming.FromContext(r.Context())
	defer func() {
		err := recover()
		if err != nil {
			create404Page(w, r, conf)
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
			debug.PrintStack()
		}
	}()
	accessLogger := loggers["access"]
	accessLogger.WithFields(logrus.Fields{
		"UA":        r.UserAgent(),
		"access_ip": r.RemoteAddr,
		"url":       r.URL.String(),
	}).Info()

	m := timing.NewMetric("f_load").Start()
	q := query.NewQueryFromGet(r)
	imageBuffer, err := content.GetImageBinary(content.GetContent(r.URL.Path, conf))
	if err != nil {
		imageBuffer = nil
		panic(imgproxyerr.New(imgproxyerr.WARNING, errors.New("can not get image data:"+err.Error())))
	}
	m.Stop()

	m = timing.NewMetric("f_decode").Start()
	img, err := imageprocessor.DecodeImage(imageBuffer)
	if err != nil {
		imageBuffer = nil
		img = nil
		panic(err)
	}
	if conf.UseMLCMYKConverter() {
		if err := img.ConvertColor(conf.MLCMYKConverterNetworkFilePath()); err != nil {
			panic(imgproxyerr.New(imgproxyerr.ERROR, err))
		}
	}
	mx, my := conf.MaxSize()
	if q.Crop() {
		img.Crop(q.Bounds().W, q.Bounds().H)
	}
	img.ResizeAndFill(q.Bounds().W, q.Bounds().H, *q.FillColor(), mx, my)
	m.Stop()

	m = timing.NewMetric("f_encode").Start()
	switch img.GetFormat() {
	case "jpeg":
		if q.UseWebP() {
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), false)
		} else if q.UseAVIF() {
			w.Header().Set("Content-Type", "image/avif")
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodeJpeg(w, img.GetImg(), q.Quality())
		}
	case "png":
		if q.UseWebP() {
			useLossless := (q.Quality() == 100)
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), useLossless)
		} else if q.UseAVIF() {
			w.Header().Set("Content-Type", "image/avif")
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodePNG(w, img.GetImg(), q.Quality())
		}
	case "gif":
		if q.UseWebP() {
			useLossless := (q.Quality() == 100)
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), useLossless)
		} else if q.UseAVIF() {
			w.Header().Set("Content-Type", "image/avif")
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodeGIF(w, img.GetImg(), q.Quality())
		}
	case "webp":
		useLossless := (q.Quality() == 100)
		err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), useLossless)
	case "avif":
		w.Header().Set("Content-Type", "image/avif")
		err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
	default:
		if q.UseWebP() {
			err = imageprocessor.EncodeWebP(w, img.GetImg(), q.Quality(), false)
		} else if q.UseAVIF() {
			w.Header().Set("Content-Type", "image/avif")
			err = imageprocessor.EncodeAVIF(w, img.GetImg(), q.Quality())
		} else {
			err = imageprocessor.EncodeJpeg(w, img.GetImg(), q.Quality())
		}
	}
	m.Stop()

	if err != nil {
		img = nil
		imageBuffer = nil
		writeDebugLog(err, conf.DebugLogPath())
		log.Println(err)

		// The following writing to the headers will be ignored if the body was wrote with some bytes.
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", "server error")
	}
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

func Initialize(conf *configure.Conf) {
	content.SetUpProviders(conf)
}
