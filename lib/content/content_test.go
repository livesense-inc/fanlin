package content

import (
	"fmt"
	"testing"

	configure "github.com/livesense-inc/fanlin/lib/conf"
)

func TestGetContent(t *testing.T) {
	t.Parallel()
	conf := configure.NewConfigure("../test/test_conf8.json")
	if conf == nil {
		t.Fatal("failed to read conf")
	}
	SetUpProviders(conf)
	cases := []struct {
		urlPath         string
		wantSourcePlace string
	}{
		{"/image.jpg", "/tmp/image.jpg"},
		{"/foo/image.jpg", "/tmp/foo/image.jpg"},
		{"/foobar/image.jpg", "/tmp/foobar/image.jpg"},
		{"/foobarbaz/image.jpg", "/tmp/foobarbaz/image.jpg"},
		{"/foobarbazgqu/image.jpg", "/tmp/foobarbazgqu/image.jpg"},
		{"/foobarbazgquu/image.jpg", "/tmp/foobarbazgquu/image.jpg"},
		{"/foobarbazgquuu/image.jpg", "/tmp/foobarbazgquuu/image.jpg"},
		{"/foobarbazgquuuu/image.jpg", "/tmp/foobarbazgquuuu/image.jpg"},
		{"/foobarbazgquuuuu/image.jpg", "/tmp/foobarbazgquuuuu/image.jpg"},
		{"/foobarbazgquuuuuu/image.jpg", "/tmp/foobarbazgquuuuuu/image.jpg"},
	}
	for n, c := range cases {
		n := n
		c := c
		t.Run(fmt.Sprintf("case-%d", n), func(t *testing.T) {
			t.Parallel()
			got := GetContent(c.urlPath, conf)
			if got == nil {
				t.Errorf("no content")
				return
			}
			if got.SourcePlace != c.wantSourcePlace {
				t.Errorf("want=%s, got=%s", c.wantSourcePlace, got.SourcePlace)
			}
		})
	}
}
