package handler

import (
	"io"
	"net/http"
	"net/url"
	"testing"

	configure "github.com/livesense-inc/fanlin/lib/conf"
	helper "github.com/livesense-inc/fanlin/lib/test"
	"github.com/sirupsen/logrus"
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

func BenchmarkMainHandler(b *testing.B) {
	c := configure.NewConfigure("../test/test_conf9.json")
	if c == nil {
		b.Fatal("Failed to build config")
	}
	Initialize(c)

	w := helper.NewNullResponseWriter()
	r := &http.Request{
		RemoteAddr: "127.0.0.1",
		Proto:      "HTTP/1.1",
		Method:     http.MethodGet,
		Host:       "127.0.0.1:3000",
		URL: &url.URL{
			Path: "/Lenna.jpg",
			RawQuery: func() string {
				q := url.Values{}
				q.Set("w", "300")
				q.Set("h", "200")
				return q.Encode()
			}(),
		},
		Header: http.Header{
			"Content-Type": []string{"application/x-www-form-urlencoded"},
		},
	}

	stdOut := logrus.New()
	stdOut.Out = io.Discard
	stdErr := logrus.New()
	stdErr.Out = io.Discard
	l := map[string]*logrus.Logger{"access": stdOut, "err": stdErr}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		MainHandler(w, r, c, l)
		if w.StatusCode() != 200 {
			b.Fatalf("want: 200, got: %d", w.StatusCode())
		}
		if w.BodySize() == 0 {
			b.Fatal("empty response body")
		}
	}
}
