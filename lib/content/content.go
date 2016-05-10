package content

import (
	"errors"

	"github.com/jobtalk/fanlin/lib/contentinfo"
	"github.com/jobtalk/fanlin/lib/error"
)

type contentType struct {
	name       string
	getContent func(*contentinfo.ContentInfo) ([]byte, error)
}

var contentTypes []contentType

// RegisterContentType registers an content type for use by GetContent.
// Name is the name of the content type, like "web" or "s3".
func RegisterContentType(name string, getContent func(*contentinfo.ContentInfo) ([]byte, error)) {
	contentTypes = append(contentTypes, contentType{
		name,
		getContent,
	})
}

// Sniff determines the contentType of c's data.
func sniff(c *contentinfo.ContentInfo) contentType {
	for _, ci := range contentTypes {
		if ci.name == c.ContentType {
			return ci
		}
	}
	return contentType{}
}

func GetContent(c *contentinfo.ContentInfo) ([]byte, error) {
	f := sniff(c)
	if f.getContent == nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("unknown content type"))
	}
	m, err := f.getContent(c)
	if err != nil {
		return nil, err
	}
	return m, nil
}
