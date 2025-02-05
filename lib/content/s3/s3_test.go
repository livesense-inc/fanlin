package s3

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/livesense-inc/fanlin/lib/content"
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

func mockS3GetFunc(config *aws.Config, bucket, key string, b *bytes.Buffer) (io.Reader, error) {
	if config == nil {
		return strings.NewReader("failed"), errors.New("config is empty")
	} else if config.Region != testRegion {
		return strings.NewReader("failed"), errors.New("Mismatch of the config region. region: " + config.Region + ", testRegion: " + testRegion)
	} else if bucket == "" {
		return strings.NewReader("failed"), errors.New("bucket is empty")
	} else if bucket != testBucket {
		return strings.NewReader("failed"), errors.New("Mismatch of the bucket. bucket: " + bucket + ", testBucket: " + testBucket)
	} else if key == "" {
		return strings.NewReader("failed"), errors.New("key is empty")
	} else if key == testKey {
		return strings.NewReader("failed"), errors.New("Mismatch of the key. key:" + key + ", testKey:" + testKey)
	}
	return strings.NewReader("success."), nil
}

func newTestContent() *content.Content {
	return &content.Content{
		SourcePlace: "s3://" + testBucket + "/" + testKey,
		SourceType:  "s3",
		Meta: map[string]interface{}{
			"region": testRegion,
		},
	}
}

func TestGetImageBinary(t *testing.T) {
	initialize()
	c := newTestContent()
	if _, err := GetImageBinary(c, new(bytes.Buffer)); err != nil {
		t.Log("normal pattern.")
		t.Fatal(err)
	}
	if _, err := GetImageBinary(nil, new(bytes.Buffer)); err == nil {
		t.Log("abnormal pattern.")
		t.Fatal("err is nil.")
	}
}
