package models

import uuid "github.com/satori/go.uuid"

type Model interface{}

type ModelImpl struct {
	ID string `json:"id"`
}

func NewModel() *ModelImpl {
	return &ModelImpl{
		ID: uuid.NewV4().String(),
	}
}
