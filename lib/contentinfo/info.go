package contentinfo

import (
	"strings"

	"github.com/google/btree"
	"github.com/jobtalk/fanlin/lib/conf"
)

type ContentInfo struct {
	ContentPlace string
	ContentType  string
	Meta         map[string]interface{}
}

type Providers struct {
	key  string
	meta interface{}
}

func (p Providers) Less(i btree.Item) bool {
	if i, ok := i.(Providers); ok {
		return -1 == strings.Compare(p.key, i.key)
	}
	return false
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
	// TODO: Review of the priorities of routing
	var bt = btree.New(4)
	for k, v := range conf.Providers() {
		bt.ReplaceOrInsert(Providers{
			k,
			v,
		})
	}

	for bt.Len() > 0 {
		root := bt.Max().(Providers)
		bt.DeleteMax()
		if strings.HasPrefix(urlPath, root.key) {
			if meta, ok := root.meta.(map[string]interface{}); ok {
				for k, info := range meta {
					switch k {
					case "src":
						src := info.(string)
						path := urlPath[len(root.key):]
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
			}
			break
		}
	}

	return &ret
}
