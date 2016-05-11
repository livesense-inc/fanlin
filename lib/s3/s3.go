package s3

import (
	"errors"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jobtalk/fanlin/lib/content"
	"github.com/jobtalk/fanlin/lib/error"
)

type S3 struct {
}

func GetSource(c *content.Content) ([]byte, error) {
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
		downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(region)}))
		_, err = downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(u.EscapedPath()),
			},
		)
		if err != nil {
			return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
		}
		bin, err := ioutil.ReadFile(file.Name())
		return bin, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	return nil, imgproxyerr.New(imgproxyerr.ERROR, errors.New("can not parse configure"))
}

func init() {
	content.RegisterContentType("s3", GetSource)
}
