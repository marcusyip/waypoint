package route

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/Sirupsen/logrus"
	"github.com/kr/pretty"
	"github.com/waypoint/waypoint/core/logger"
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
	logger    *logrus.Logger
	taskRepo  repos.RouteTaskRepository
	mapClient maps.Client
}

func (mgr *RouteManagerImpl) getLogger(method string) *logrus.Entry {
	return mgr.logger.WithFields(logrus.Fields{"manager": "RouteManager", "method": method})
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
	log := mgr.getLogger("RunTask").WithField("route_task_id", routeTaskID)
	m, err := mgr.taskRepo.Get(routeTaskID)
	if err != nil {
		log.WithField("err", err).Info("Failed to get route task data")
		return err
	}
	task := m.(*models.RouteTask)

	err = mgr.validateTask(task)
	if err != nil {
		err2 := mgr.saveError(task, err.Error())
		if err2 != nil {
			log = log.WithField("err2", err2)
		}
		return err
	}

	r := mgr.getDirectionsRequest(task)
	resp, _, err := mgr.mapClient.Directions(context.Background(), r)
	if err != nil {
		err2 := mgr.saveError(task, err.Error())
		if err2 != nil {
			log = log.WithField("err2", err2)
		}
		log.WithField("err", err).Info("Response with error")
		return err
	}
	pretty.Print(resp)
	log.Info("Successfully calculate route")
	return mgr.saveResult(task, resp)
}

func (mgr *RouteManagerImpl) validateTask(task *models.RouteTask) error {
	if len(task.Route) < 2 {
		return errors.New("Missing origin or destination")
	}
	return nil
}

func (mgr *RouteManagerImpl) getDirectionsRequest(task *models.RouteTask) *gmaps.DirectionsRequest {
	origin := task.Route[0]
	destination := task.Route[len(task.Route)-1]
	waypoints := make([]string, 0, len(task.Route)-2)
	for _, point := range task.Route[1 : len(task.Route)-1] {
		waypoints = append(waypoints, fmt.Sprintf("%s,%s", point[0], point[1]))
	}
	return &gmaps.DirectionsRequest{
		Mode:        gmaps.TravelModeDriving,
		Origin:      fmt.Sprintf("%s,%s", origin[0], origin[1]),
		Destination: fmt.Sprintf("%s,%s", destination[0], destination[1]),
		Waypoints:   waypoints,
	}
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
	for _, step := range route.Legs[0].Steps {
		result.TotalDistance += step.Distance.Meters
		result.TotalTime += step.Duration.Seconds()
		endLocation := step.EndLocation
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
		logger:    logger.GetLogger(),
		taskRepo:  repos.GetRouteTaskRepository(),
		mapClient: maps.GetClient(),
	}
}
