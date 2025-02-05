package local

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/livesense-inc/fanlin/lib/content"
)

func GetImageBinary(c *content.Content, buf *bytes.Buffer) (io.Reader, error) {
	f, err := os.Open(path.Clean(c.SourcePlace))
	if err != nil {
		return nil, fmt.Errorf("failed to open a file: %s: %w", c.SourcePlace, err)
	}
	defer f.Close()
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}
	return buf, nil
}

func init() {
	content.RegisterContentType("local", GetImageBinary)
}
