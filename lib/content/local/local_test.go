package local

import (
	"io"
	"testing"

	"github.com/livesense-inc/fanlin/lib/content"
)

func TestGetImageBinary(t *testing.T) {
	c := content.Content{
		SourcePlace: "../../test/img/Lenna.jpg",
	}
	if r, err := GetImageBinary(&c, []byte{}); err != nil {
		t.Fatal(err)
	} else {
		if b, err := io.ReadAll(r); err != nil {
			t.Fatal(err)
		} else {
			if len(b) == 0 {
				t.Error("something was wrong: zero byte")
			}
		}
	}
}
