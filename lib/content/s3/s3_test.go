package s3

import (
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/livesense-inc/fanlin/lib/content"
	"io"
	"strings"
)

var (
	SetS3GetFunc    = setS3GetFunc
	CreateAwsConfig = createAwsConfig
	testBucket      = "testBucket"
	testRegion      = "ap-northeast-1"
	testKey         = "test/test.jpg"
)

func initialize() {
	SetS3GetFunc(mockS3GetFunc)
	testBucket = "testBucket"
	testRegion = "ap-northeast-1"
	testKey = "test/test.jpg"
}

func mockS3GetFunc(config *aws.Config, bucket, key string, file *os.File) (io.Reader, error) {
	if config == nil {
		return strings.NewReader("failed"), errors.New("config is empty")
	} else if aws.StringValue(config.Region) != testRegion {
		return strings.NewReader("failed"), errors.New("Mismatch of the config region. region: " + aws.StringValue(config.Region) + ", testRegion: " + testRegion)
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

func TestCreateAwsConfig(t *testing.T) {
	region := "ap-northeast-1"
	meta := map[string]interface{}{}
	awsConfig := CreateAwsConfig(region, meta)

	if aws.StringValue(awsConfig.Region) != region {
		t.Fatalf("Mismatch of the region. '%s' expected, '%s' got.", region, aws.StringValue(awsConfig.Region))
	}

	if awsConfig.Credentials != nil {
		t.Fatalf("Unexpected Credentials got. nil edpected, '%v' got.", awsConfig.Credentials)
	}

	meta["use_env_credential"] = true
	awsConfig = CreateAwsConfig(region, meta)
	if awsConfig.Credentials == nil {
		t.Fatal("NIL Credentials got.")
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
