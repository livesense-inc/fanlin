package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"

	configure "github.com/livesense-inc/fanlin/lib/conf"
	"github.com/livesense-inc/fanlin/lib/handler"
	"github.com/livesense-inc/fanlin/lib/logger"
	servertiming "github.com/mitchellh/go-server-timing"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/netutil"
)

var confList = []string{
	"/etc/fanlin.json",
	"/etc/fanlin.cnf",
	"/etc/fanlin.conf",
	"/usr/local/etc/fanlin.json",
	"/usr/local/etc/fanlin.cnf",
	"/usr/local/etc/fanlin.conf",
	"/usr/local/fanlin/fanlin.json",
	"/usr/local/fanlin.json",
	"/usr/local/fanlin.conf",
	"./fanlin.json",
	"./fanlin.cnf",
	"./fanlin.conf",
	"./conf.json",
	"~/.fanlin.json",
	"~/.fanlin.cnf",
}

var (
	buildVersion string
	buildHash    string
	buildDate    string
	goversion    string
)

var (
	vOption bool
)

func showVersion() {
	fmt.Println()
	if buildVersion != "" {
		fmt.Println("build version: ", buildVersion)
	} else {
		fmt.Println("build version: ", buildHash)
	}
	fmt.Println("build date: ", buildDate)
	fmt.Println("GO version: ", goversion)
}

func main() {
	conf := func() *configure.Conf {
		for _, confName := range confList {
			if conf, path := configure.NewConfigure(confName); conf != nil {
				fmt.Println("read configure. :", path)
				return conf
			}
		}
		return nil
	}()
	if conf == nil {
		log.Fatal("Can not read configure.")
	}

	notFoundImagePath := flag.String("nfi", conf.NotFoundImagePath(), "not found image path")
	errorLogPath := flag.String("err", conf.ErrorLogPath(), "error log path")
	accessLogPath := flag.String("log", conf.AccessLogPath(), "access log path")
	port := flag.Int("p", conf.Port(), "port")
	port = flag.Int("port", *port, "port")
	maxProcess := flag.Int("cpu", conf.MaxProcess(), "max process.")
	debug := flag.Bool("debug", false, "debug mode.")
	flag.BoolVar(&vOption, "v", false, "version")
	flag.Parse()

	if vOption {
		showVersion()
		os.Exit(128)
	}

	conf.Set("404_img_path", *notFoundImagePath)
	conf.Set("error_log_path", *errorLogPath)
	conf.Set("access_log_path", *accessLogPath)
	conf.Set("port", *port)
	conf.Set("max_process", *maxProcess)

	loggers := map[string]*logrus.Logger{
		"err":    logger.NewLogger(conf.ErrorLogPath()),
		"access": logger.NewLogger(conf.AccessLogPath()),
	}

	if *debug {
		log.Print(conf)
	}

	http.DefaultClient.Timeout = conf.BackendRequestTimeout()
	runtime.GOMAXPROCS(conf.MaxProcess())

	fn := func(w http.ResponseWriter, r *http.Request) {
		handler.MainHandler(w, r, conf, loggers)
	}
	var h http.Handler = http.HandlerFunc(fn)
	if conf.UseServerTiming() {
		h = servertiming.Middleware(h, nil)
	}
	http.Handle("/", h)
	http.HandleFunc("/healthCheck", handler.HealthCheckHandler)

	if conf.EnableMetricsEndpoint() {
		metricsHandler := handler.MakeMetricsHandler(conf, log.New(os.Stderr, "", log.LstdFlags))
		http.HandleFunc("/metrics", metricsHandler.ServeHTTP)
	}

	addr := fmt.Sprintf(":%d", conf.Port())
	if n := conf.MaxClients(); n > 0 {
		ln, err := net.Listen("tcp", addr)
		if err != nil {
			log.Fatal(err)
		}
		lln := netutil.LimitListener(ln, n)
		http.Serve(lln, nil)
	} else {
		http.ListenAndServe(addr, nil)
	}
}
