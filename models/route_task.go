package models

var (
	RouteTaskStatusPending = "pending"
	RouteTaskStatusSuccess = "success"
	RouteTaskStatusError   = "error"
)

type RouteTask struct {
	ModelImpl
	Status string          `json:"status"`
	Route  [][]string      `json:"route"`
	Result RouteTaskResult `json:"result"`
	Reason string          `json:"reason"`
}

type RouteTaskResult struct {
	Path          [][]string `json:"path"`
	TotalDistance int        `json:"total_distance"`
	TotalTime     float64    `json:"total_time"`
}

func NewRouteTask() *RouteTask {
	return &RouteTask{
		ModelImpl: *NewModel(),
		Status:    RouteTaskStatusPending,
	}
}
