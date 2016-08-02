package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/livesense-inc/fanlin/lib/conf"
	"github.com/livesense-inc/fanlin/lib/content"
	_ "github.com/livesense-inc/fanlin/lib/content/s3"
	_ "github.com/livesense-inc/fanlin/lib/content/web"
	"github.com/livesense-inc/fanlin/lib/error"
	"github.com/livesense-inc/fanlin/lib/image"
	"github.com/livesense-inc/fanlin/lib/query"
	_ "github.com/livesense-inc/fanlin/plugin"
	"github.com/sirupsen/logrus"
)

var devNull, _ = os.Open("/dev/null")

func create404Page(w http.ResponseWriter, r *http.Request, conf *configure.Conf) {
	q := query.NewQueryFromGet(r)

	maxW, maxH := conf.MaxSize()
	bin, err := imageprocessor.Set404Image(conf.NotFoundImagePath(), q.Bounds().W, q.Bounds().H, *q.FillColor(), maxW, maxH)

	w.WriteHeader(404)
	if err != nil {
		fmt.Fprintf(w, "%s", "404 Not found.")
	} else {
		fmt.Fprintf(w, "%s", bin)
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

func MainHandler(w http.ResponseWriter, r *http.Request, conf *configure.Conf, loggers map[string]logrus.Logger) {
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
						break
					case imgproxyerr.ERROR:
						writeDebugLog(err, conf.DebugLogPath())
						errLogger.Error(err)
						break
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

	q := query.NewQueryFromGet(r)
	imageBuffer, err := content.GetImageBinary(content.GetContent(r.URL.Path, conf))

	if err != nil {
		imageBuffer = nil
		panic(imgproxyerr.New(imgproxyerr.WARNING, errors.New("can not get image data:"+err.Error())))
	}

	img, err := imageprocessor.DecodeImage(imageBuffer)
	if err != nil {
		imageBuffer = nil
		img = nil
		panic(err)
	}
	mx, my := conf.MaxSize()
	img.ResizeAndFill(q.Bounds().W, q.Bounds().H, *q.FillColor(), mx, my)

	imageBuffer, err = imageprocessor.EncodeJpeg(img.GetImg())
	if err != nil {
		img = nil
		imageBuffer = nil
		panic(err)
	}

	fmt.Fprintf(w, "%s", imageBuffer)
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", "")
}
