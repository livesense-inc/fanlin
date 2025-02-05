package web

import (
	"fmt"
	"net/http"
	"runtime"

	"bytes"
	"io"

	"github.com/livesense-inc/fanlin/lib/content"
	imgproxyerr "github.com/livesense-inc/fanlin/lib/error"
)

var ua = fmt.Sprintf("Mozilla/5.0 (fanlin; arch: %s; OS: %s; Go version: %s) Go language Client/1.1 (KHTML, like Gecko) Version/1.0 fanlin", runtime.GOARCH, runtime.GOOS, runtime.Version())

var httpClient = Client{
	Http: new(RealWebClient),
}

type RealWebClient struct {
}

type WebClient interface {
	Get(string, []byte) (io.Reader, error)
}

type Client struct {
	Http WebClient
}

func (r *RealWebClient) Get(url string, b []byte) (io.Reader, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	req.Header.Set("User-Agent", ua)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	defer resp.Body.Close()

	if isErrorCode(resp.StatusCode) {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, fmt.Errorf("received error status code(%d)", resp.StatusCode))
	}

	buffer := bytes.NewBuffer(b)
	if _, err := io.Copy(buffer, resp.Body); err != nil {
		return nil, err
	}

	return buffer, nil
}

func isErrorCode(status int) bool {
	switch status / 100 {
	case 4, 5:
		return true
	default:
		return false
	}
}

func GetImageBinary(c *content.Content, b []byte) (io.Reader, error) {
	return httpClient.Http.Get(c.SourcePlace, b)
}

func setHttpClient(c Client) {
	httpClient = c
}

func init() {
	content.RegisterContentType("web", GetImageBinary)
}
