package handler

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bass3m/swarmeus/config"
	"github.com/bass3m/swarmeus/scan"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Index page\n")
}

func Status(c config.Config) func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Status page\n")
		targets, err := scan.GetTargets(c)
		if err != nil {
			log.Errorf("Failed to targets:  %v ", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(targets)
	}
}

func SetupRoutes(router *httprouter.Router, c config.Config, routePrefix string) {
	router.GET("/", Index)
	router.GET(routePrefix+"/status", Status(c))
}
