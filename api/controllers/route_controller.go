package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"github.com/waypoint/waypoint/core/api"
	"github.com/waypoint/waypoint/core/logger"
	"github.com/waypoint/waypoint/entities"
	routemgr "github.com/waypoint/waypoint/managers/route"
	"github.com/waypoint/waypoint/queue"
	repos "github.com/waypoint/waypoint/repositories"
)

type RouteController struct {
	api.Controller
	logger   *logrus.Logger
	routeMgr routemgr.RouteManager
	routeEnt entities.RouteEntity
}

func (c *RouteController) getLogger(method string) *logrus.Entry {
	return c.logger.WithFields(logrus.Fields{"controller": "RouteController", "method": method})
}

func (c *RouteController) Show(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := ps.ByName("token")
	log := c.getLogger("Show").WithField("id", id)
	task, err := c.routeMgr.GetByID(id)
	if err != nil {
		log.WithField("err", err).Info("Failed to get route result by ID")
		if err == repos.ErrNotFound {
			c.Error(rw, 404, 10001, "Not found")
			return
		}
		c.Error(rw, 400, 99999, "Unknown error")
		return
	}
	log.WithField("err", err).Info("Successfully get route result by ID")
	c.JSON(rw, 200, c.routeEnt.New(task))
}

func (c *RouteController) Create(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	log := c.getLogger("Create")
	var route [][]string
	err := json.NewDecoder(r.Body).Decode(&route)
	if err != nil {
		log.WithField("err", err).Info("Unable to unmarshal JSON")
		c.Error(rw, 400, 10001, err.Error())
		return
	}
	task, err := c.routeMgr.CreateAsyncTask(queue.GetServer(), route)
	if err != nil {
		log.WithField("err", err).Info("Failed to create task")
		c.Error(rw, 400, 99999, "Unknown error")
		return
	}
	resp := map[string]string{
		"token": task.ID,
	}
	log.WithField("task_id", task.ID).Info("Successfully create task")
	c.JSON(rw, 201, resp)
}

func GetRouteController() *RouteController {
	return &RouteController{
		logger:   logger.GetLogger(),
		routeMgr: routemgr.GetRouteManager(),
		routeEnt: entities.GetRouteEntity(),
	}
}
