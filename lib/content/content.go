package content

import (
	"net/url"
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

var providers []provider

const DEFAULT_PRIORITY float64 = 10.0

func init() {
	providers = nil
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

func GetContent(urlPath string, conf *configure.Conf) *Content {
	if urlPath == "" {
		return nil
	}
	if conf == nil {
		return nil
	}

	if providers == nil {
		providers = getProviders(conf)
	}

	return getContent(urlPath, providers)
}
