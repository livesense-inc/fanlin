package handler

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	configure "github.com/livesense-inc/fanlin/lib/conf"
	"github.com/sirupsen/logrus"
)

// NullResponseWriter is a struct
type NullResponseWriter struct {
	h          http.Header
	b          io.Writer
	StatusCode int
}

// NewNullResponseWriter returns a instance of NullResponseWriter
func NewNullResponseWriter() *NullResponseWriter {
	return &NullResponseWriter{h: http.Header{}, b: io.Discard}
}

// Header returns headers
func (w *NullResponseWriter) Header() http.Header {
	return w.h
}

// Write writes bytes to response writer
func (w *NullResponseWriter) Write(p []byte) (int, error) {
	return w.b.Write(p)
}

// WriteHeader updates the status code
func (w *NullResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

func BenchmarkMainHandler(b *testing.B) {
	w := NewNullResponseWriter()

	r := &http.Request{
		RemoteAddr: "127.0.0.1",
		Proto:      "HTTP/1.1",
		Method:     http.MethodGet,
		Host:       "127.0.0.1:8080",
		URL: &url.URL{
			Path: "/test/master/granada-jp0802/858.jpg",
		},
		Header: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
		},
	}

	c := configure.NewConfigure("test.json")
	if c == nil {
		b.Fatal("Failed to build config")
	}

	stdOut := logrus.New()
	stdOut.Out = io.Discard
	stdErr := logrus.New()
	stdErr.Out = io.Discard
	l := map[string]logrus.Logger{"access": *stdOut, "err": *stdErr}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MainHandler(w, r, c, l)
		if w.StatusCode != 200 {
			b.Fatalf("want: 200, got: %d", w.StatusCode)
		}
	}
}
