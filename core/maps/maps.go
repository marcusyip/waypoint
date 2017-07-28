package maps

import (
	"github.com/waypoint/waypoint/core/config"
	gmaps "googlemaps.github.io/maps"
)

var (
	client *gmaps.Client
)

func GetClient() *gmaps.Client {
	return client
}

func Init() {
	conf := config.GetConfig().GoogleAPI
	if len(conf.APIKey) == 0 {
		panic("Missing Google Maps API key")
	}
	var err error
	client, err = gmaps.NewClient(gmaps.WithAPIKey(conf.APIKey))
	if err != nil {
		panic(err)
	}
}
