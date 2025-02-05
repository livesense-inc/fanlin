package s3

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/livesense-inc/fanlin/lib/content"
	imgproxyerr "github.com/livesense-inc/fanlin/lib/error"
	"golang.org/x/text/unicode/norm"
)

var s3GetSourceFunc = getS3ImageBinary

// Test dedicated function
func setS3GetFunc(f func(cfg *aws.Config, bucket, key string, b *bytes.Buffer) (io.Reader, error)) {
	s3GetSourceFunc = f
}

func GetImageBinary(c *content.Content, b *bytes.Buffer) (io.Reader, error) {
	if c == nil {
		return nil, errors.New("content is nil")
	}
	s3url := c.SourcePlace
	u, err := url.Parse(s3url)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("can not parse s3 url"))
	}

	bucket := u.Host

	if region, ok := c.Meta["region"].(string); ok {
		path, err := url.QueryUnescape(u.EscapedPath())
		if err != nil {
			return nil, err
		}
		if form, ok := c.Meta["norm_form"].(string); ok {
			path, err = NormalizePath(path, form)
			if err != nil {
				return nil, err
			}
		}
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return nil, err
		}
		return s3GetSourceFunc(&cfg, bucket, path, b)
	}
	return nil, imgproxyerr.New(imgproxyerr.ERROR, errors.New("can not parse configure"))
}

func NormalizePath(path string, form string) (string, error) {
	switch form {
	case "nfd":
		return norm.NFD.String(path), nil
	case "nfc":
		return norm.NFC.String(path), nil
	case "nfkc":
		return norm.NFKC.String(path), nil
	case "nfkd":
		return norm.NFKD.String(path), nil
	}
	return "", imgproxyerr.New(imgproxyerr.WARNING, errors.New("invalid normalization form("+form+")"))
}

func getS3ImageBinary(cfg *aws.Config, bucket, key string, buf *bytes.Buffer) (io.Reader, error) {
	cli := s3.NewFromConfig(*cfg)
	input := s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	output, err := cli.GetObject(context.TODO(), &input)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}
	defer output.Body.Close()
	if _, err := io.Copy(buf, output.Body); err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}
	return buf, nil
}

func init() {
	content.RegisterContentType("s3", GetImageBinary)
}
