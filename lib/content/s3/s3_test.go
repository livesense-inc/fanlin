package s3

import (
	"errors"
	"os"
	"testing"

	"github.com/livesense-inc/fanlin/lib/content"
	"io"
	"strings"
)

var (
	SetS3GetFunc = setS3GetFunc
	testBucket   = "testBucket"
	testRegion   = "ap-northeast-1"
	testKey      = "test/test.jpg"
)

func initialize() {
	SetS3GetFunc(mockS3GetFunc)
	testBucket = "testBucket"
	testRegion = "ap-northeast-1"
	testKey = "test/test.jpg"
}

func mockS3GetFunc(region, bucket, key string, file *os.File) (io.Reader, error) {
	if region == "" {
		return strings.NewReader("failed"), errors.New("region is empty")
	} else if region != testRegion {
		return strings.NewReader("failed"), errors.New("Mismatch of the region. region: " + region + ", testRegion: " + testRegion)
	} else if bucket == "" {
		return strings.NewReader("failed"), errors.New("bucket is empty")
	} else if bucket != testBucket {
		return strings.NewReader("failed"), errors.New("Mismatch of the bucket. bucket: " + bucket + ", testBucket: " + testBucket)
	} else if key == "" {
		return strings.NewReader("failed"), errors.New("key is empty")
	} else if key == testKey {
		return strings.NewReader("failed"), errors.New("Mismatch of the key. key:" + key + ", testKey:" + testKey)
	} else if file == nil {
		return strings.NewReader("failed"), errors.New("file is nil")
	}
	return strings.NewReader("success."), nil
}

func newTestContent() *content.Content {
	return &content.Content{
		"s3://" + testBucket + "/" + testKey,
		"s3",
		map[string]interface{}{
			"region": testRegion,
		},
	}
}

func TestGetImageBinary(t *testing.T) {
	initialize()
	c := newTestContent()
	if _, err := GetImageBinary(c); err != nil {
		t.Log("normal pattern.")
		t.Fatal(err)
	}
	if _, err := GetImageBinary(nil); err == nil {
		t.Log("abnormal pattern.")
		t.Fatal("err is nil.")
	}
}
