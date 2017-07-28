package entities

import "github.com/waypoint/waypoint/models"

type Route struct {
	Status        string     `json:"status"`
	Path          [][]string `json:"path,omitempty"`
	TotalDistance int        `json:"total_distance,omitempty"`
	TotalTime     float64    `json:"total_time,omitempty"`
	Error         string     `json:"error,omitempty"`
}

type RouteEntity struct{}

func (e *RouteEntity) New(task *models.RouteTask) *Route {
	r := &Route{}
	switch task.Status {
	case models.RouteTaskStatusPending:
		r.Status = "in progress"
	case models.RouteTaskStatusError:
		r.Status = "failure"
		r.Error = task.Reason
	case models.RouteTaskStatusSuccess:
		r.Status = "success"
		r.Path = task.Result.Path
		r.TotalDistance = task.Result.TotalDistance
		r.TotalTime = task.Result.TotalTime
	}
	return r
}

func GetRouteEntity() RouteEntity {
	return RouteEntity{}
}
