package route

import (
	routemgr "github.com/waypoint/waypoint/managers/route"
)

type RouteTask struct {
	routeMgr routemgr.RouteManager
}

func (t *RouteTask) Run(routeTaskID string) error {
	return t.routeMgr.RunTask(routeTaskID)
}

func GetRouteTask() *RouteTask {
	return &RouteTask{
		routeMgr: routemgr.GetRouteManager(),
	}
}
