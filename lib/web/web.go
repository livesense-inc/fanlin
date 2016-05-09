package web

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/jobtalk/fanlin/lib/content"
	"github.com/jobtalk/fanlin/lib/contentinfo"
	"github.com/jobtalk/fanlin/lib/error"
)

type Web struct {
	ua string
}

var client = http.Client{
	Transport: &http.Transport{MaxIdleConnsPerHost: 64},
	Timeout:   time.Duration(10) * time.Second,
}

func New(ua string) *Web {
	return &Web{ua}
}

func isErrorCode(status int) bool {
	switch status / 100 {
	case 4, 5:
		return true
	default:
		return false
	}
}

func GetContent(c *contentinfo.ContentInfo) ([]byte, error) {
	req, err := http.NewRequest("GET", c.ContentPlace, nil)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
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
	content.RegisterContentType("web", GetContent)
}
