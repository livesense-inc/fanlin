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

func convertInterfaceToMap(i interface{}) map[string]interface{} {
	if ret, ok := i.(map[string]interface{}); ok {
		return ret
	}
	return map[string]interface{}(nil)
}

func GetContentInfo(urlPath string, conf *configure.Conf) *ContentInfo {
	var ret ContentInfo
	ret.Meta = map[string]interface{}{}
	if urlPath == "" {
		return nil
	}
	if conf == nil {
		return nil
	}
	for _, buffer := range conf.Providers() {
		provider := convertInterfaceToMap(buffer)
		for key, meta := range provider {
			if strings.HasPrefix(urlPath, key) {
				for k, info := range convertInterfaceToMap(meta) {
					switch k {
					case "src":
						src := info.(string)
						path := urlPath[len(key):]
						if !strings.HasPrefix(path, "/") {
							path = "/" + path
						}
						ret.ContentPlace = src + path
					case "type":
						ret.ContentType = info.(string)
					default:
						ret.Meta[k] = info
					}
				}
				return &ret
			}
		}
	}
	return nil
}
