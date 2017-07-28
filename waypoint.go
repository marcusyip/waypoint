package main

import (
	"os"

	"github.com/waypoint/waypoint/api"
	"github.com/waypoint/waypoint/core/config"
	"github.com/waypoint/waypoint/core/maps"
	"github.com/waypoint/waypoint/core/redis"
	"github.com/waypoint/waypoint/queue"

	cli "gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "wayponit"
	app.Usage = "Routing service"
	app.Action = func(c *cli.Context) error {
		start()
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "Run task for calculating shortest driving path and estimated driving time",
			Action: func(c *cli.Context) error {
				startWorker()
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func start() {
	loadConfig()

	// redis
	pool := redis.Init()
	defer pool.Close()
	// start queue server
	_, err := queue.StartServer()
	if err != nil {
		panic(err)
	}
	// API server
	server := api.NewServer()
	server.Start()
}

func startWorker() {
	loadConfig()

	// redis
	pool := redis.Init()
	defer pool.Close()
	// google maps client
	maps.Init()

	server, err := queue.StartServer()
	if err != nil {
		panic(err)
	}
	worker := server.NewWorker("waypoint_worker")
	if err := worker.Launch(); err != nil {
		panic(err)
	}
}

func loadConfig() {
	if err := config.Load("/etc/waypoint/config.json"); err == nil {
		return
	}
	if err := config.Load("./etc/waypoint/config.json"); err == nil {
		return
	}
	panic("[main] Failed to load configuration file")
}
