package web

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	"github.com/jobtalk/fanlin/lib/content"
	"github.com/jobtalk/fanlin/lib/error"
)

var ua = fmt.Sprintf("Mozilla/5.0 (fanlin; arch: %s; OS: %s; Go version: %s) Go language Client/1.1 (KHTML, like Gecko) Version/1.0 fanlin", runtime.GOARCH, runtime.GOOS, runtime.Version())

type Web struct {
}

var client = http.Client{
	Transport: &http.Transport{MaxIdleConnsPerHost: 64},
	Timeout:   time.Duration(10) * time.Second,
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
	req, err := http.NewRequest("GET", c.SourcePlace, nil)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	req.Header.Set("User-Agent", ua)
	res, err := client.Do(req)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	} else if isErrorCode(res.StatusCode) {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("Image can not get"))
	}

	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}
	return data, nil
}

func init() {
	content.RegisterContentType("web", GetSource)
}
