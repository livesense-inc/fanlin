package content

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	iradix "github.com/hashicorp/go-immutable-radix"
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
	providers []provider
	router    *iradix.Tree
	useRouter bool
)

const DEFAULT_PRIORITY float64 = 10.0

func init() {
	providers = nil
	router = nil
	useRouter = false
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

func makeRouter(conf *configure.Conf) *iradix.Tree {
	// https://pkg.go.dev/github.com/hashicorp/go-immutable-radix
	router := iradix.New()
	for _, p := range conf.Providers() {
		for alias, meta := range convertInterfaceToMap(p) {
			prefix := strings.TrimSuffix(alias, "/")
			if !strings.HasPrefix(prefix, "/") {
				prefix = fmt.Sprintf("/%s", prefix)
			}
			m := convertInterfaceToMap(meta)
			router, _, _ = router.Insert([]byte(prefix), provider{alias, m})
		}
	}
	return router
}

func getContent(urlPath string, p []provider, r *iradix.Tree, useRouter bool) *Content {
	if urlPath == "/" || urlPath == "" {
		return nil
	}
	targetProvider := searchProvider(urlPath, p, r, useRouter)
	if targetProvider == nil {
		return nil
	}
	var ret Content
	ret.Meta = map[string]interface{}{}
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

func searchProvider(urlPath string, p []provider, r *iradix.Tree, useRouter bool) *provider {
	if useRouter {
		return searchProviderByRouter(urlPath, r)
	}

	return searchProviderFromSortedList(urlPath, p)
}

func searchProviderByRouter(urlPath string, r *iradix.Tree) *provider {
	if _, value, ok := r.Root().LongestPrefix([]byte(urlPath)); ok {
		if provider, ok := value.(provider); ok {
			return &provider
		}
	}
	return nil
}

func searchProviderFromSortedList(urlPath string, p []provider) *provider {
	index := -1
	for i, v := range p {
		if strings.HasPrefix(urlPath, v.alias) {
			index = i
			break
		}
	}
	if index == -1 {
		return nil
	}
	return &p[index]
}

func convertInterfaceToMap(i interface{}) map[string]interface{} {
	if ret, ok := i.(map[string]interface{}); ok {
		return ret
	}
	return map[string]interface{}(nil)
}

func SetUpProviders(conf *configure.Conf) {
	providers = getProviders(conf)
	router = makeRouter(conf)
	priorities := make(map[float64]struct{}, len(providers))
	for _, provider := range providers {
		if p, ok := provider.meta["priority"].(float64); ok {
			priorities[p] = struct{}{}
		}
	}
	useRouter = len(priorities) <= 1
}

func GetContent(urlPath string, conf *configure.Conf) *Content {
	if urlPath == "" {
		return nil
	}
	if conf == nil {
		return nil
	}
	if providers == nil {
		return nil
	}
	if router == nil {
		return nil
	}
	return getContent(urlPath, providers, router, useRouter)
}
