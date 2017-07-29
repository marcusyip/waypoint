package route

import (
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/waypoint/waypoint/core/logger"
	"github.com/waypoint/waypoint/core/maps"
	"github.com/waypoint/waypoint/mocks"
	"github.com/waypoint/waypoint/models"
	gmaps "googlemaps.github.io/maps"
)

var _ = Describe("Route Manager", func() {
	var (
		mockCtrl          *gomock.Controller
		routeManager      *RouteManagerImpl
		mockRouteTaskRepo *mocks.MockRouteTaskRepository
		mockMapClient     *maps.MockClient
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		mockRouteTaskRepo = mocks.NewMockRouteTaskRepository(mockCtrl)
		mockMapClient = maps.NewMockClient(mockCtrl)
		routeManager = &RouteManagerImpl{
			logger:    logger.GetLogger(),
			taskRepo:  mockRouteTaskRepo,
			mapClient: mockMapClient,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("Run task", func() {
		It("should make request to map API and update route task result", func() {
			mockRouteTask := models.NewRouteTask()
			mockRouteTask.Route = [][]string{
				[]string{"10", "10"},
				[]string{"10.2", "10.2"},
			}
			mockGMapRoutes := []gmaps.Route{
				{
					Legs: []*gmaps.Leg{
						{
							Steps: []*gmaps.Step{
								{
									StartLocation: gmaps.LatLng{
										Lat: 10,
										Lng: 10,
									},
									EndLocation: gmaps.LatLng{
										Lat: 10.1,
										Lng: 10.1,
									},
									Distance: gmaps.Distance{
										Meters: 10,
									},
									Duration: time.Duration(20 * time.Second),
								},
								{
									StartLocation: gmaps.LatLng{
										Lat: 10.1,
										Lng: 10.1,
									},
									EndLocation: gmaps.LatLng{
										Lat: 10.2,
										Lng: 10.2,
									},
									Distance: gmaps.Distance{
										Meters: 20,
									},
									Duration: time.Duration(30 * time.Second),
								},
							},
						},
					},
				},
			}
			mockRouteTaskRepo.EXPECT().Get(mockRouteTask.ID).Times(1).Return(mockRouteTask, nil)
			mockMapClient.EXPECT().Directions(gomock.Any(), gomock.Any()).Times(1).Return(mockGMapRoutes, []gmaps.GeocodedWaypoint{}, nil)
			mockRouteTaskRepo.EXPECT().Set(gomock.Any()).Times(1).Do(func(m *models.RouteTask) {
				Expect(m.Status).To(Equal(models.RouteTaskStatusSuccess))
				Expect(m.Result.TotalDistance).To(Equal(30))
				Expect(m.Result.TotalTime).To(Equal(float64(50)))
				Expect(len(m.Result.Path)).To(Equal(3))
			})

			err := routeManager.RunTask(mockRouteTask.ID)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
