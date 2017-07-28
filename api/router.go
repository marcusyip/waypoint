package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	ctrls "github.com/waypoint/waypoint/api/controllers"
)

func WelcomeHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprint(rw, "Welcome to XForum")
}

func withRouter(router *httprouter.Router) {
	fmt.Printf("[router] initializing router\n")

	routeCtrl := ctrls.GetRouteController()
	router.GET("/route/:token", routeCtrl.Show)
	router.POST("/route", routeCtrl.Create)
}
