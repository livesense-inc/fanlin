package content

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/jobtalk/fanlin/lib/client"
	"github.com/jobtalk/fanlin/lib/conf"
)

type Content struct {
	data []byte
}

func getSourceURL(urlPath string, conf *configure.Conf) (ret string) {
	if urlPath == "" {
		return ""
	}
	if conf == nil {
		return ""
	}
	for k, v := range conf.Externals() {
		searchWord := "/" + k + "/"
		if strings.HasPrefix(urlPath, searchWord) {
			return v + urlPath[len(searchWord):]
		}
	}
	return ""
}

func (c *Content) getExternalContent(urlPath string, conf *configure.Conf) ([]byte, bool) {
	bin, err := client.HttpImageGetter(getSourceURL(urlPath, conf), conf)
	return bin, err == nil
}

func (c *Content) getLocalContent(urlPath string, conf *configure.Conf) ([]byte, bool) {
	localPath, err := filepath.Abs(conf.LocalImagePath())
	if err != nil {
		return nil, err == nil
	}
	requestPath, err := filepath.Abs(localPath + urlPath)
	if err != nil {
		return nil, err == nil
	}
	if strings.Contains(requestPath, localPath) {
		bin, err := ioutil.ReadFile(requestPath)
		return bin, err == nil
	}
	return nil, false
}

func (c *Content) getS3Content(urlPath string, conf *configure.Conf) ([]byte, bool) {
	if bin, err := client.S3ImageGetter(urlPath, conf); err == nil {
		return bin, true
	}
	return nil, false
}

func GetContent(urlPath string, conf *configure.Conf) ([]byte, bool) {
	c := &Content{}
	if ret, ok := c.getExternalContent(urlPath, conf); ok {
		return ret, ok
	} else if ret, ok := c.getS3Content(urlPath, conf); ok {
		return ret, ok
	} else if ret, ok := c.getLocalContent(urlPath, conf); ok {
		return ret, ok
	}
	return nil, false
}
