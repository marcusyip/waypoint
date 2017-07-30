# Run API server, worker and redis in docker
```
docker-compose up
```

server should be running on 3000 port
```
http://{docker_ip}:3000
```

# Structure

waypoint.go` - app main

`api` - API server, router, controller

`managers` - core business login

`queue` - task queue service

`tasks` - Task entry point

`core` - core modules

`repositories` - data access layer

`models` - data models 

`entities` - view models, entity builder

`etc` - config
