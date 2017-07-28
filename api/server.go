package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/urfave/negroni"
	"github.com/waypoint/waypoint/core/config"
)

type Server struct{}

func (s *Server) Start() {
	conf := config.GetConfig().API

	n := negroni.New(negroni.HandlerFunc(loggingHandler))
	router := httprouter.New()
	withRouter(router)
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%d", conf.Port))
}

func NewServer() *Server {
	return &Server{}
}

func loggingHandler(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now().UTC()
	logrus.WithFields(logrus.Fields{"method": r.Method, "path": r.URL.Path}).Info("HTTP request")
	next(rw, r)
	logrus.WithFields(logrus.Fields{"method": r.Method, "path": r.URL.Path, "resp_time": time.Since(start)}).Info("HTTP request")
}
