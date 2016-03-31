package client

import (
	"testing"

	"github.com/jobtalk/fanlin/lib/conf"
)

var (
	IsErrorCode = isErrorCode
	conf        = configure.NewConfigure("../test/test_conf.json")
)

func TestHttpImageGetter(t *testing.T) {
	data, err := HttpImageGetter("https://google.co.jp/", conf)
	if err != nil {
		t.Error(err)
	}
	if len(data) <= 0 {
		t.Fatalf("Faild. Can not get body.")
	}
}

func TestIsErrorCode(t *testing.T) {
	t.Log("ステータスコードが4xx, 5xxの判定テスト")

	if IsErrorCode(200) {
		t.Error(200, IsErrorCode(200))
	}
	if !IsErrorCode(400) {
		t.Error(400, !IsErrorCode(400))
	}
	if !IsErrorCode(499) {
		t.Error(499, !IsErrorCode(499))
	}
	if !IsErrorCode(500) {
		t.Error(500, !IsErrorCode(500))
	}
	if !IsErrorCode(599) {
		t.Error(599, !IsErrorCode(599))
	}
	if IsErrorCode(600) {
		t.Error(600, IsErrorCode(600))
	}
}
