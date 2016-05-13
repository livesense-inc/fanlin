package web

import (
	"errors"
	"testing"

	"github.com/jobtalk/fanlin/lib/content"
)

var SetHttpClient = setHttpClient
var (
	IsErrorCode = isErrorCode
	targetURL   = "https://google.co.jp"
	bin         = []byte("It works.")
)

type MockWebClient struct {
}

func getTestClient() *Client {
	c := new(Client)
	c.Http = new(MockWebClient)
	return c
}

func (mwc *MockWebClient) Get(url string) ([]byte, error) {
	if url != targetURL {
		return nil, errors.New("not match url. url: " + url + ", targetURL: " + targetURL)
	}
	return bin, nil
}

func TestIsErrorCode(t *testing.T) {
	if IsErrorCode(200) {
		t.Fatal(200, IsErrorCode(200))
	}
	if IsErrorCode(203) {
		t.Fatal(203, IsErrorCode(203))
	}
	if !IsErrorCode(404) {
		t.Fatal(404, IsErrorCode(404))
	}
	if !IsErrorCode(500) {
		t.Fatal(500, IsErrorCode(500))
	}
}

func TestGetImageBinary(t *testing.T) {
	SetHttpClient(*getTestClient())
	c := &content.Content{
		SourcePlace: targetURL,
	}
	result, err := GetImageBinary(c)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != "It works." {
		t.Fatal(string(result))
	}
}
