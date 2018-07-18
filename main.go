package main

import (
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/bass3m/swarmeus/config"
	"github.com/bass3m/swarmeus/handler"
	"github.com/bass3m/swarmeus/scan"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	var (
		app = kingpin.New(filepath.Base(os.Args[0]), "swarmeus")

		listenAddress = app.Flag("web.listen-address", "Address to listen on for the web interface and API.").Default(":9723").String()
		routePrefix   = app.Flag("web.route-prefix", "Prefix for the internal routes of web endpoints.").Default("").String()
		configPath    = app.Flag("cfg.path", "Path to YAML configuration file.").Default("/etc/swarmeus/swarmeus.yml").String()
	)
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *routePrefix == "/" {
		*routePrefix = ""
	}
	if *routePrefix != "" {
		*routePrefix = "/" + strings.Trim(*routePrefix, "/")
	}

	log.Infoln("Starting swarmeus")
	log.Debugf("Prefix path is '%s'", *routePrefix)

	flags := map[string]string{}
	for _, f := range app.Model().Flags {
		flags[f.Name] = f.Value.String()
	}

	c, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugln("Config:", c)

	// start scanning for containers with metrics
	cancel := make(chan struct{})
	go scan.Scan(c, cancel)

	router := httprouter.New()
	handler.SetupRoutes(router, *routePrefix)

	l, err := net.Listen("tcp", *listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	go interruptHandler(l, cancel)
	err = (&http.Server{Addr: *listenAddress, Handler: router}).Serve(l)
	log.Errorln("swarmeus HTTP server stopped:", err)
}

func interruptHandler(l net.Listener, cancel chan<- struct{}) {
	var e struct{}
	notifier := make(chan os.Signal)
	signal.Notify(notifier, os.Interrupt, syscall.SIGTERM)
	<-notifier
	log.Info("swarmeus Received SIGINT/SIGTERM; exiting ...")
	l.Close()
	// send cancel event
	cancel <- e
}
