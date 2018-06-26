package s3

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/url"
	"os"

	"golang.org/x/text/unicode/norm"

	"io"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/livesense-inc/fanlin/lib/content"
	"github.com/livesense-inc/fanlin/lib/error"
)

var s3GetSourceFunc = getS3ImageBinary

// Test dedicated function
func setS3GetFunc(f func(region, bucket, key string, file *os.File) (io.Reader, error)) {
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

	file, err := ioutil.TempFile(os.TempDir(), "s3_img")
	defer func() {
		os.Remove(file.Name())
		file.Close()
	}()
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
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
			}
		}
		return s3GetSourceFunc(region, bucket, path, file)
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

func getS3ImageBinary(region, bucket, key string, file *os.File) (io.Reader, error) {
	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(region)}))
	_, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		},
	)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}

	ret := new(bytes.Buffer)
	if _, err := io.Copy(ret, file); err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}
	return ret, imgproxyerr.New(imgproxyerr.ERROR, err)
}

func init() {
	content.RegisterContentType("s3", GetImageBinary)
}
