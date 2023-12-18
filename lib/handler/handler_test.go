package handler

import (
	"net/http"
	"net/url"
	"testing"

	configure "github.com/livesense-inc/fanlin/lib/conf"
	helper "github.com/livesense-inc/fanlin/lib/test"
)

func TestMakeMetricsHandler(t *testing.T) {
	r := &http.Request{
		RemoteAddr: "127.0.0.1",
		Proto:      "HTTP/1.1",
		Method:     http.MethodGet,
		Host:       "127.0.0.1:8080",
		URL: &url.URL{
			Path: "/metrics",
		},
	}
	w := helper.NewNullResponseWriter()
	h := MakeMetricsHandler(&configure.Conf{}, helper.NullLogger())
	h.ServeHTTP(w, r)
	if w.StatusCode() != http.StatusOK {
		t.Errorf("want=%d, got=%d", http.StatusOK, w.StatusCode())
	}
}
