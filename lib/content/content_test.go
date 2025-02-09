package content

import (
	"testing"

	configure "github.com/livesense-inc/fanlin/lib/conf"
)

func TestGetContent(t *testing.T) {
	conf := configure.NewConfigure("../test/test_conf8.json")
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

	for _, flag := range []bool{false, true} {
		useRouter = flag
		for n, c := range cases {
			got := GetContent(c.urlPath, conf)
			if got == nil {
				t.Errorf("useRouter: %v, case: %d, no content", useRouter, n)
				continue
			}
			if got.SourcePlace != c.wantSourcePlace {
				t.Errorf(
					"useRouter: %v, case: %d, want=%s, got=%s",
					useRouter,
					n,
					c.wantSourcePlace,
					got.SourcePlace,
				)
			}
		}
	}

}

func BenchmarkGetContentFromSortedList(b *testing.B) {
	conf := configure.NewConfigure("../test/test_conf7.json")
	SetUpProviders(conf)
	useRouter = false
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got := GetContent("/foo", conf); got == nil {
			b.Fatalf("no content")
		}
	}
}

func BenchmarkGetContentByRouter(b *testing.B) {
	conf := configure.NewConfigure("../test/test_conf7.json")
	SetUpProviders(conf)
	useRouter = true
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if got := GetContent("/foobarbazgquuuuuu", conf); got == nil {
			b.Fatalf("no content")
		}
	}
}
