package repositories

import "github.com/waypoint/waypoint/models"

type RouteTaskRepository interface {
	Get(id string) (models.Model, error)
	Set(m *models.RouteTask) error
}

type RouteTaskRepositoryImpl struct{}

func (r *RouteTaskRepositoryImpl) Get(id string) (models.Model, error) {
	return get(r, r.getKey(id))
}

func (r *RouteTaskRepositoryImpl) Set(m *models.RouteTask) error {
	return set(r, m)
}

func (r *RouteTaskRepositoryImpl) getKey(id string) string {
	return "route_task:" + id
}

func (r *RouteTaskRepositoryImpl) getModelKey(m models.Model) string {
	t := m.(*models.RouteTask)
	return r.getKey(t.ID)
}

func (r *RouteTaskRepositoryImpl) getModel() models.Model {
	return &models.RouteTask{}
}

func GetRouteTaskRepository() RouteTaskRepository {
	return &RouteTaskRepositoryImpl{}
}
