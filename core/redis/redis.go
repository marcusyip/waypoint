package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/waypoint/waypoint/core/config"
)

var (
	pool *redis.Pool
)

func GetPool() *redis.Pool {
	return pool
}

func Init() *redis.Pool {
	conf := config.GetConfig().Redis
	pool = redis.NewPool(func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", conf.Address)
		if err != nil {
			return nil, err
		}
		return c, nil
	}, 10)
	return pool
}
