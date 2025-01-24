package local

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/livesense-inc/fanlin/lib/content"
)

func GetImageBinary(c *content.Content) (io.Reader, error) {
	log.Print(c.SourcePlace)
	f, err := os.Open(c.SourcePlace)
	if err != nil {
		return nil, fmt.Errorf("Failed to open a file: %s: %w", c.SourcePlace, err)
	}
	defer f.Close()
	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}
	return &b, nil
}

func init() {
	content.RegisterContentType("local", GetImageBinary)
}
