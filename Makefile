.PHONY: mockgen

mockgen:
	@echo Generate mocks folder

	mockgen -destination=mocks/repositories.go -package=mocks github.com/waypoint/waypoint/repositories RouteTaskRepository
	mockgen -destination=core/maps/mocks.go -package=maps github.com/waypoint/waypoint/core/maps Client
