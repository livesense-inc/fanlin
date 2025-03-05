package s3

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/livesense-inc/fanlin/lib/content"
	imgproxyerr "github.com/livesense-inc/fanlin/lib/error"
	"golang.org/x/text/unicode/norm"
)

var s3GetSourceFunc = getS3ImageBinary

// Test dedicated function
func setS3GetFunc(f func(cli *s3.Client, bucket, key string) (io.Reader, error)) {
	s3GetSourceFunc = f
}

func GetImageBinary(c *content.Content) (io.Reader, error) {
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
		useMock := false
		if v, ok := c.Meta["use_mock"]; ok {
			if b, ok := v.(bool); ok {
				useMock = b
			}
		}
		cli, err := makeClient(region, useMock)
		if err != nil {
			return nil, err
		}
		return s3GetSourceFunc(cli, bucket, path)
	}
	return nil, imgproxyerr.New(imgproxyerr.ERROR, errors.New("can not parse configure"))
}

func makeClient(region string, useMock bool) (*s3.Client, error) {
	if !useMock {
		cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
		if err != nil {
			return nil, err
		}
		return s3.NewFromConfig(
			cfg,
			func(o *s3.Options) {
				o.DisableLogOutputChecksumValidationSkipped = true
			},
		), nil
	}

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				"AAAAAAAAAAAAAAAAAAAA",
				"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
				"",
			),
		),
	)
	if err != nil {
		return nil, err
	}
	return s3.NewFromConfig(
		cfg,
		s3.WithEndpointResolverV2(
			s3.EndpointResolverV2(
				s3.NewDefaultEndpointResolverV2(),
			),
		),
		func(o *s3.Options) {
			o.BaseEndpoint = aws.String("http://localhost:4567")
			o.UsePathStyle = true
			o.DisableLogOutputChecksumValidationSkipped = true
		},
	), nil
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

func getS3ImageBinary(cli *s3.Client, bucket, key string) (io.Reader, error) {
	key = strings.TrimPrefix(key, "/")
	downloader := s3manager.NewDownloader(cli)
	buf := s3manager.NewWriteAtBuffer([]byte{})
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err := downloader.Download(context.TODO(), buf, input)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, fmt.Errorf("bucket=%s, key=%s: %w", bucket, key, err))
	}
	return bytes.NewReader(buf.Bytes()), nil
}

func init() {
	content.RegisterContentType("s3", GetImageBinary)
}
