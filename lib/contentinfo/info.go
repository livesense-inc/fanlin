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

func GetContentInfo(urlPath string, conf *configure.Conf) *ContentInfo {
	var ret ContentInfo
	if urlPath == "" {
		return nil
	}
	if conf == nil {
		return nil
	}
	for k, v := range conf.Providers() {
		searchWord := "/" + k + "/"
		if strings.HasPrefix(urlPath, searchWord) {
			if infos, ok := v.(map[string]interface{}); ok {
				for k, info := range infos {
					switch k {
					case "type":
						buf, ok := info.(string)
						if !ok {
							return nil
						}
						ret.ContentType = buf
					case "src":
						buf, ok := info.(string)
						if !ok {
							return nil
						}
						ret.ContentPlace = buf + urlPath[len(searchWord):]
					default:
						ret.Meta[k] = info
					}
				}
			}
		}
	}
	return &ret
}
