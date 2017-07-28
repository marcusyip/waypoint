package route

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/waypoint/waypoint/core/maps"
	"github.com/waypoint/waypoint/models"
	repos "github.com/waypoint/waypoint/repositories"
	gmaps "googlemaps.github.io/maps"
)

type RouteManager interface {
	GetByID(id string) (*models.RouteTask, error)
	CreateAsyncTask(queueServer *machinery.Server, route [][]string) (*models.RouteTask, error)
	RunTask(routeTaskID string) error
}

type RouteManagerImpl struct {
	taskRepo repos.RouteTaskRepository
}

func (mgr *RouteManagerImpl) GetByID(id string) (*models.RouteTask, error) {
	m, err := mgr.taskRepo.Get(id)
	if err != nil {
		return nil, err
	}
	task := m.(*models.RouteTask)
	return task, nil
}

func (mgr *RouteManagerImpl) CreateAsyncTask(queueServer *machinery.Server, route [][]string) (*models.RouteTask, error) {
	task := models.NewRouteTask()
	task.Route = route
	err := mgr.taskRepo.Set(task)
	if err != nil {
		return nil, err
	}
	signature := &tasks.Signature{
		UUID: task.ID,
		Name: "route",
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: task.ID,
			},
		},
	}
	_, err = queueServer.SendTask(signature)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (mgr *RouteManagerImpl) RunTask(routeTaskID string) error {
	m, err := mgr.taskRepo.Get(routeTaskID)
	if err != nil {
		return err
	}
	task := m.(*models.RouteTask)
	if len(task.Route) < 2 {
		err := mgr.saveError(task, "Missing origin or destination")
		if err != nil {
			return err
		}
		return errors.New("Missing origin or destination")
	}
	origin := task.Route[0]
	destination := task.Route[len(task.Route)-1]
	waypoints := make([]string, 0, len(task.Route)-2)
	fmt.Printf("this %+v", task)
	for _, point := range task.Route[1 : len(task.Route)-1] {
		waypoints = append(waypoints, fmt.Sprintf("%s,%s", point[0], point[1]))
	}
	c := maps.GetClient()
	r := &gmaps.DirectionsRequest{
		Mode:        gmaps.TravelModeDriving,
		Origin:      fmt.Sprintf("%s,%s", origin[0], origin[1]),
		Destination: fmt.Sprintf("%s,%s", destination[0], destination[1]),
		Waypoints:   waypoints,
	}
	resp, _, err := c.Directions(context.Background(), r)
	if err != nil {
		return err
	}
	return mgr.saveResult(task, resp)
}

func (mgr *RouteManagerImpl) saveResult(task *models.RouteTask, routes []gmaps.Route) error {
	task.Status = models.RouteTaskStatusSuccess
	task.Result = *mgr.getResult(routes)
	return mgr.taskRepo.Set(task)
}

func (mgr *RouteManagerImpl) saveError(task *models.RouteTask, reason string) error {
	task.Status = models.RouteTaskStatusError
	task.Reason = reason
	return mgr.taskRepo.Set(task)
}

func (mgr *RouteManagerImpl) getResult(routes []gmaps.Route) *models.RouteTaskResult {
	result := &models.RouteTaskResult{}
	route := routes[0]
	result.Path = make([][]string, 0, len(route.Legs))
	startLocation := route.Legs[0].StartLocation
	result.Path = append(result.Path, latLng(startLocation.Lat, startLocation.Lng))
	for _, leg := range route.Legs {
		result.TotalDistance += leg.Distance.Meters
		result.TotalTime += leg.Duration.Seconds()
		endLocation := leg.EndLocation
		result.Path = append(result.Path, latLng(endLocation.Lat, endLocation.Lng))
	}
	return result
}

func latLng(lat float64, lng float64) []string {
	return []string{
		strconv.FormatFloat(lat, 'f', -1, 64),
		strconv.FormatFloat(lng, 'f', -1, 64),
	}
}

func GetRouteManager() RouteManager {
	return &RouteManagerImpl{
		taskRepo: repos.GetRouteTaskRepository(),
	}
}
