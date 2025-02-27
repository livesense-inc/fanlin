package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

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
			if conf := configure.NewConfigure(confName); conf != nil {
				if p, err := filepath.Abs(confName); err != nil {
					fmt.Println("read configure. :", confName)
				} else {
					fmt.Println("read configure. :", p)
				}
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

	if err := handler.Prepare(conf); err != nil {
		log.Fatal(err)
	}

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

	if err := runServer(conf); err != nil {
		log.Fatal(err)
	}

	loggers["err"].Print("shut down")
	os.Exit(0)
}

func runServer(conf *configure.Conf) error {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", conf.Port()))
	if err != nil {
		return err
	}
	if conf.MaxClients() > 0 {
		listener = netutil.LimitListener(listener, conf.MaxClients())
	}
	defer listener.Close()

	server := &http.Server{}
	if conf.ServerTimeout() > 0*time.Second {
		server.ReadTimeout = conf.ServerTimeout()
		server.WriteTimeout = conf.ServerTimeout()
	}
	if conf.ServerIdleTimeout() > 0*time.Second {
		server.IdleTimeout = conf.ServerIdleTimeout()
	}

	c := make(chan os.Signal, 1)
	defer close(c)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer signal.Stop(c)

	go func(s *http.Server, c <-chan os.Signal) {
		if _, ok := <-c; ok {
			s.SetKeepAlivesEnabled(false)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			_ = s.Shutdown(ctx)
		}
	}(server, c)

	if err := server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
