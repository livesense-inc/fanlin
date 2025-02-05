package content

import (
	"bytes"
	"errors"

	"io"

	imgproxyerr "github.com/livesense-inc/fanlin/lib/error"
)

type source struct {
	name           string
	getImageBinary func(*Content, *bytes.Buffer) (io.Reader, error)
}

var sources []source

// RegisterContentType registers an content type for use by GetContent.
// Name is the name of the content type, like "web" or "s3".
func RegisterContentType(name string, getImageBinary func(*Content, *bytes.Buffer) (io.Reader, error)) {
	sources = append(sources, source{
		name,
		getImageBinary,
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

func GetImageBinary(c *Content, b *bytes.Buffer) (io.Reader, error) {
	f := sniff(c)
	if f.getImageBinary == nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("unknown content type"))
	}
	m, err := f.getImageBinary(c, b)
	if err != nil {
		return nil, err
	}
	return m, nil
}
