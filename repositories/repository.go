package repositories

import (
	"encoding/json"
	"errors"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/waypoint/waypoint/core/redis"
	"github.com/waypoint/waypoint/models"
)

var (
	ErrNotFound = errors.New("repositories: not found")
)

type Repository interface {
	getModel() models.Model
	getModelKey(m models.Model) string
}

func get(r Repository, key string) (models.Model, error) {
	c := getConn()
	c.Send("GET", key)
	c.Flush()
	v, err := c.Receive()
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, ErrNotFound
	}
	m := r.getModel()
	b := v.([]byte)
	err = json.Unmarshal(b, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func set(r Repository, m models.Model) error {
	c := getConn()
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	c.Send("SET", r.getModelKey(m), b)
	c.Flush()
	_, err = c.Receive()
	if err != nil {
		return err
	}
	return nil
}

func getConn() redigo.Conn {
	return redis.GetPool().Get()
}
