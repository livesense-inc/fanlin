package test

import (
	"io"
	"log"
	"net/http"
	"time"
)

// for test
const (
	Timeout = 5 * time.Second
)

// NullResponseWriter is a struct
type NullResponseWriter struct {
	h     http.Header
	b     io.Writer
	bSize int
	sc    int
}

// NewNullResponseWriter returns a instance of NullResponseWriter
func NewNullResponseWriter() *NullResponseWriter {
	return &NullResponseWriter{h: http.Header{}, b: io.Discard, sc: http.StatusOK}
}

// Header returns headers
func (w *NullResponseWriter) Header() http.Header {
	return w.h
}

// Write writes bytes to response writer
func (w *NullResponseWriter) Write(p []byte) (int, error) {
	w.bSize += len(p)
	return w.b.Write(p)
}

// WriteHeader updates the status code
func (w *NullResponseWriter) WriteHeader(statusCode int) {
	w.sc = statusCode
}

// StatusCode returns a status code
func (w *NullResponseWriter) StatusCode() int {
	return w.sc
}

// BodySize returns a size of the body
func (w *NullResponseWriter) BodySize() int {
	return w.bSize
}

// NullLogger returns a instance of log.Logger
func NullLogger() *log.Logger {
	return log.New(io.Discard, "", 0)
}
