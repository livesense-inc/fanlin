package client

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/jobtalk/fanlin/lib/conf"
	"github.com/jobtalk/fanlin/lib/error"
)

func S3ImageGetter(urlPath string, conf *configure.Conf) ([]byte, error) {
	path := strings.TrimLeft(urlPath, "/")
	if path == "" {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("not exist"))
	}

	file, err := ioutil.TempFile(os.TempDir(), "s3_img")
	defer func() {
		os.Remove(file.Name())
		file.Close()
	}()
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}

	downloader := s3manager.NewDownloader(session.New(&aws.Config{Region: aws.String(conf.S3Region())}))
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(conf.S3BucketName()),
			Key:    aws.String(path),
		},
	)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}

	bin, err := ioutil.ReadFile(file.Name())
	return bin, imgproxyerr.New(imgproxyerr.ERROR, err)
}
