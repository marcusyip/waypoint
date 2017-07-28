package queue

import (
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/waypoint/waypoint/core/config"
	routetask "github.com/waypoint/waypoint/tasks/route"
)

var (
	server *machinery.Server
)

func StartServer() (*machinery.Server, error) {
	conf := config.GetConfig().Machinery
	var err error
	server, err = machinery.NewServer(&conf)

	tasks := map[string]interface{}{
		"route": routetask.GetRouteTask().Run,
	}
	err = server.RegisterTasks(tasks)
	if err != nil {
		return nil, err
	}
	return server, nil
}

func GetServer() *machinery.Server {
	return server
}
