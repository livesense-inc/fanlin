package contentinfo

import (
	"strings"

	"github.com/jobtalk/fanlin/lib/conf"
)

type ContentInfo struct {
	ContentPlace string
	ContentType  string
	Meta         map[string]interface{}
}

type provider struct {
	alias string
	meta  interface{}
}

func getProviders(c *configure.Conf) []provider {
	ret := []provider{}
	for _, p := range c.Providers() {
		for alias, meta := range convertInterfaceToMap(p) {
			ret = append(ret, provider{alias, meta})
		}
	}
	return ret
}

func getContentInfo(urlPath string, p []provider) *ContentInfo {
	var ret ContentInfo
	ret.Meta = map[string]interface{}{}
	index := serachProviderIndex(urlPath, p)
	if index < 0 {
		return nil
	}
	targetProvider := p[index]
	for k, v := range convertInterfaceToMap(targetProvider.meta) {
		switch k {
		case "src":
			src := v.(string)
			path := urlPath[len(targetProvider.alias):]
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}
			ret.ContentPlace = src + path
		case "type":
			ret.ContentType = v.(string)
		default:
			ret.Meta[k] = v
		}
	}
	return &ret
}

func serachProviderIndex(urlPath string, p []provider) int {
	for i, v := range p {
		if strings.Contains(urlPath, v.alias) {
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

func GetContentInfo(urlPath string, conf *configure.Conf) *ContentInfo {
	if urlPath == "" {
		return nil
	}
	if conf == nil {
		return nil
	}
	providers := getProviders(conf)

	return getContentInfo(urlPath, providers)
}
