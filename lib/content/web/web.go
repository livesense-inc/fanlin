package web

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/jobtalk/fanlin/lib/content"
	"github.com/jobtalk/fanlin/lib/error"
)

var ua = fmt.Sprintf("Mozilla/5.0 (fanlin; arch: %s; OS: %s; Go version: %s) Go language Client/1.1 (KHTML, like Gecko) Version/1.0 fanlin", runtime.GOARCH, runtime.GOOS, runtime.Version())

var client = http.Client{
	Transport: &http.Transport{MaxIdleConnsPerHost: 64},
	Timeout:   time.Duration(10) * time.Second,
}

var httpClient = Client{
	Http: new(RealWebClient),
}

type RealWebClient struct {
}

type WebClient interface {
	Get(string) ([]byte, error)
}

type Client struct {
	Http WebClient
}

func (r *RealWebClient) Get(url string) ([]byte, error) {
	var body []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	req.Header.Set("User-Agent", ua)

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	return body, nil
}

func isErrorCode(status int) bool {
	switch status / 100 {
	case 4, 5:
		return true
	default:
		return false
	}
}

func GetSource(c *content.Content) ([]byte, error) {
	return httpClient.Http.Get(c.SourcePlace)
}

func setHttpClient(c Client) {
	httpClient = c
}

func init() {
	content.RegisterContentType("web", GetSource)
}
