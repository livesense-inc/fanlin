package configure

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Conf struct {
	c map[string]interface{}
}

func (c *Conf) UseMLCMYKConverter() bool {
	b, ok := c.c["use_ml_cmyk_converter"]
	if !ok {
		return false
	}
	r, ok := b.(bool)
	if !ok {
		panic("'use_ml_cmyk_converter' parameter is incorrect")
	}
	return r
}

func (c *Conf) MLCMYKConverterNetworkFilePath() string {
	i, ok := c.c["ml_cmyk_converter_network_file_path"]
	if !ok {
		return ""
	}
	s, ok := i.(string)
	if !ok {
		panic("'ml_cmyk_converter_network_file_path' parameter is incorrect")
	}
	return s
}

func (c *Conf) UseServerTiming() bool {
	b, ok := c.c["use_server_timing"]
	if !ok {
		return false
	}
	r, ok := b.(bool)
	if !ok {
		panic("'use_server_timing' parameter is incorrect")
	}
	return r
}

func (c *Conf) EnableMetricsEndpoint() bool {
	b, ok := c.c["enable_metrics_endpoint"]
	if !ok {
		return false
	}
	r, ok := b.(bool)
	if !ok {
		panic("'enable_metrics_endpoint' parameter is incorrect")
	}
	return r
}

func (c *Conf) MaxClients() int {
	b, ok := c.c["max_clients"]
	if !ok {
		return 0
	}
	n, ok := b.(float64)
	if !ok {
		panic("'max_clients' parameter is incorrect")
	}
	return int(n)
}

func (c *Conf) Set(k string, v interface{}) {
	if v == nil {
		return
	}
	c.c[k] = v
}

func (c *Conf) Get(k string) interface{} {
	return c.c[k]
}

func (c *Conf) UA() string {
	ua := c.c["user_agent"]
	if ua == nil {
		ua = fmt.Sprintf("Mozilla/5.0 (fanlin; arch: %s; OS: %s; Go version: %s) Go language Client/1.1 (KHTML, like Gecko) Version/1.0 fanlin", runtime.GOARCH, runtime.GOOS, runtime.Version())
	}
	return ua.(string)
}

func (c *Conf) Providers() []interface{} {
	if providers, ok := c.c["providers"].([]interface{}); ok {
		return providers
	}
	return nil
}

func (c *Conf) NotFoundImagePath() string {
	path := c.c["404_img_path"]
	if path == nil {
		path = "./"
	}
	return path.(string)
}

func (c *Conf) ErrorLogPath() string {
	path := c.c["error_log_path"]
	if path == nil {
		path = "./error.log"
	}
	return path.(string)
}

func (c *Conf) AccessLogPath() string {
	path := c.c["access_log_path"]
	if path == nil {
		path = "./access.log"
	}
	return path.(string)
}

func (c *Conf) DebugLogPath() string {
	path := c.c["debug_log_path"]
	if path == nil {
		path = "/dev/null"
	}
	return path.(string)
}

func (c *Conf) BackendRequestTimeout() time.Duration {
	tstr, ok := c.c["backend_request_timeout"]
	if !ok {
		return 10 * time.Second
	}

	t, ok := tstr.(string)
	if !ok {
		return 10 * time.Second
	}

	d, err := time.ParseDuration(t)
	if err != nil {
		panic(err)
	}
	return d
}

func (c *Conf) Port() int {
	port := c.c["port"]
	return convInterfaceToInt(port, 8080)
}

func (c *Conf) MaxSize() (uint, uint) {
	width := c.c["max_width"]
	height := c.c["max_height"]
	w := convInterfaceToInt(width, 1000)
	h := convInterfaceToInt(height, 1000)

	return uint(w), uint(h)
}

func (c *Conf) MaxProcess() int {
	maxProcess := c.c["max_process"]
	return convInterfaceToInt(maxProcess, runtime.NumCPU())
}

func convInterfaceToInt(v interface{}, exception int) int {
	if n, ok := v.(float64); ok {
		return int(n)
	} else if n, ok := v.(int); ok {
		return n
	} else if n, ok := v.(uint); ok {
		return int(n)
	}
	return exception
}

func (c *Conf) getIncludePath() []string {
	if includes, ok := c.c["include"].([]interface{}); ok {
		ret := make([]string, len(includes))
		for i, s := range includes {
			if path, ok := s.(string); ok {
				ret[i] = path
			}
		}
		return ret
	}
	return nil
}

func toAbsPath(mainConfPath string, includeConfPath string) string {
	var includeAbs string
	mainConfAbs, err := filepath.Abs(mainConfPath)
	if err != nil {
		return ""
	}
	if filepath.IsAbs(includeConfPath) {
		return includeConfPath
	}
	includeAbs, err = filepath.Abs(mainConfAbs + "/" + includeConfPath)
	if err != nil {
		return ""
	}
	return includeAbs
}

func (c *Conf) includeConfigure(mainConfPath string, pathList []string) {
	parentConfPath := func() string {
		d := strings.Split(mainConfPath, "/")
		d = d[:len(d)-1]
		return strings.Join(d, "/")
	}()
	for _, path := range pathList {
		abs := toAbsPath(parentConfPath, path)
		if abs == "" {
			continue
		}
		include := NewConfigure(abs)
		if include == nil {
			continue
		}

		for k, v := range include.c {
			c.c[k] = v
		}
	}
}

func NewConfigure(confPath string) *Conf {
	var conf map[string]interface{}
	bin, err := os.ReadFile(confPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = json.Unmarshal(bin, &conf)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	p, _ := filepath.Abs(confPath)
	fmt.Println("read configure. :", p)
	c := Conf{
		c: conf,
	}
	includes := c.getIncludePath()
	c.includeConfigure(confPath, includes)

	delete(c.c, "include")
	return &c
}
