package content

import (
	"errors"

	"github.com/jobtalk/fanlin/lib/error"
)

type source struct {
	name       string
	getContent func(*Content) ([]byte, error)
}

var sources []source

// RegisterContentType registers an content type for use by GetContent.
// Name is the name of the content type, like "web" or "s3".
func RegisterContentType(name string, getContent func(*Content) ([]byte, error)) {
	sources = append(sources, source{
		name,
		getContent,
	})
}

// Sniff determines the contentType of c's data.
func sniff(c *Content) source {
	for _, ci := range sources {
		if ci.name == c.SourceType {
			return ci
		}
	}
	return source{}
}

func GetSource(c *Content) ([]byte, error) {
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
