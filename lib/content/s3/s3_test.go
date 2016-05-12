package s3

import (
	"errors"
	"os"
	"testing"

	"github.com/jobtalk/fanlin/lib/content"
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

func mockS3GetFunc(region, bucket, key string, file *os.File) ([]byte, error) {
	if region == "" {
		return []byte("failed."), errors.New("region is empty")
	} else if region != testRegion {
		return []byte("failed."), errors.New("Mismatch of the region. region: " + region + ", testRegion: " + testRegion)
	} else if bucket == "" {
		return []byte("failed."), errors.New("bucket is empty")
	} else if bucket != testBucket {
		return []byte("failed."), errors.New("Mismatch of the bucket. bucket: " + bucket + ", testBucket: " + testBucket)
	} else if key == "" {
		return []byte("failed."), errors.New("key is empty")
	} else if key == testKey {
		return []byte("failed."), errors.New("Mismatch of the key. key:" + key + ", testKey:" + testKey)
	} else if file == nil {
		return []byte("failed."), errors.New("file is nil")
	}
	return []byte("success."), nil
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

func TestGetSource(t *testing.T) {
	initialize()
	c := newTestContent()
	if _, err := GetSource(c); err != nil {
		t.Log("normal pattern.")
		t.Fatal(err)
	}
	if _, err := GetSource(nil); err == nil {
		t.Log("abnormal pattern.")
		t.Fatal("err is nil.")
	}
}
