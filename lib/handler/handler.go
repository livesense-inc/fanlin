package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/jobtalk/fanlin/lib/conf"
	"github.com/jobtalk/fanlin/lib/content"
	"github.com/jobtalk/fanlin/lib/contentinfo"
	"github.com/jobtalk/fanlin/lib/error"
	"github.com/jobtalk/fanlin/lib/image"
	"github.com/jobtalk/fanlin/lib/query"
	_ "github.com/jobtalk/fanlin/lib/s3"
	_ "github.com/jobtalk/fanlin/lib/web"
	"github.com/sirupsen/logrus"
)

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
						errLogger.Warn(err)
						break
					case imgproxyerr.ERROR:
						errLogger.Error(err)
						break
					default:
						errLogger.Error(err)
					}
				} else {
					errLogger.Error(err)
				}

			} else {
				log.Println(err)
			}
			fmt.Fprintf(w, "%s", "")
		}
	}()

	accessLogger := loggers["access"]
	accessLogger.WithFields(logrus.Fields{
		"UA":        r.UserAgent(),
		"access_ip": r.RemoteAddr,
		"url":       r.URL.String(),
	}).Info()

	q := query.NewQueryFromGet(r)
	info := contentinfo.GetContentInfo(r.URL.Path, conf)
	imageBuffer, err := content.GetContent(info)

	if err != nil {
		imageBuffer = nil
		panic(imgproxyerr.New(imgproxyerr.WARNING, errors.New("can not get image data")))
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
