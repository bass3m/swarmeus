package handler

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Index page\n")
}

func Status(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Status page\n")
}

func SetupRoutes(router *httprouter.Router, routePrefix string) {
	router.GET("/", Index)
	router.GET(routePrefix+"/status", Status)
}
