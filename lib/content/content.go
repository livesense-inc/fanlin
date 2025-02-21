package content

import (
	"bytes"
	"io"
	"net/url"
	"os"
	"sort"
	"strings"

	configure "github.com/livesense-inc/fanlin/lib/conf"
)

type Content struct {
	SourcePlace string
	SourceType  string
	Meta        map[string]interface{}
}

type provider struct {
	alias string
	meta  map[string]interface{}
}

var (
	providers      []provider
	noContentImage []byte
)

const DEFAULT_PRIORITY float64 = 10.0

func init() {
	providers = nil
	noContentImage = nil
}

func getProviders(c *configure.Conf) []provider {
	ret := make([]provider, 0, len(c.Providers()))
	for _, p := range c.Providers() {
		for alias, meta := range convertInterfaceToMap(p) {
			m := convertInterfaceToMap(meta)
			if _, ok := m["priority"]; !ok {
				m["priority"] = DEFAULT_PRIORITY
			}

			ret = append(ret, provider{alias, m})
		}
	}

	sort.Slice(ret, func(i, j int) bool {
		return ret[i].meta["priority"].(float64) < ret[j].meta["priority"].(float64)
	})

	return ret
}

func getContent(urlPath string, p []provider) *Content {
	if urlPath == "/" || urlPath == "" {
		return nil
	}
	var ret Content
	ret.Meta = map[string]interface{}{}
	index := serachProviderIndex(urlPath, p)
	if index < 0 {
		return nil
	}
	targetProvider := p[index]
	for k, v := range targetProvider.meta {
		switch k {
		case "src":
			src := v.(string)
			path := urlPath[len(targetProvider.alias):]
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			ret.SourcePlace, _ = url.QueryUnescape(src + path)
		case "type":
			ret.SourceType = v.(string)
		default:
			ret.Meta[k] = v
		}
	}
	return &ret
}

func serachProviderIndex(urlPath string, p []provider) int {
	for i, v := range p {
		if strings.HasPrefix(urlPath, v.alias) {
			return i
		}
	}
	return -1
}

func convertInterfaceToMap(i interface{}) map[string]interface{} {
	if ret, ok := i.(map[string]interface{}); ok {
		return ret
	}
	return map[string]interface{}(nil)
}

func loadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var b bytes.Buffer
	if _, err := io.Copy(&b, f); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func SetupNoContentImage(conf *configure.Conf) (err error) {
	if conf == nil {
		return
	}

	noContentImage, err = loadFile(conf.NotFoundImagePath())
	return
}

func SetUpProviders(conf *configure.Conf) {
	if conf == nil {
		return
	}

	providers = getProviders(conf)
}

func GetNoContentImage() io.Reader {
	return bytes.NewReader(noContentImage)
}

func GetContent(urlPath string, _conf *configure.Conf) *Content {
	if urlPath == "" {
		return nil
	}

	return getContent(urlPath, providers)
}
