# Run API server, worker and redis in docker
```
docker-compose up
```

# Structure

`xforum.go` - app main
`api` - API server, router, controller
`managers` - core business login
`queue` - task queue service
`tasks` - Task entry point
`core` - core modules
`repositories` - data access layer
`models` - data models 
`entities` - view models, entity builder
`etc` - config
