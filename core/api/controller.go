package api

import (
	"net/http"

	"github.com/unrolled/render"
)

type Controller struct{}

func (c *Controller) JSON(rw http.ResponseWriter, status int, v interface{}) {
	r := render.New()
	r.JSON(rw, status, v)
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Controller) Error(rw http.ResponseWriter, status, errCode int, errMessage string) {
	r := render.New()
	r.JSON(rw, status, ErrorResponse{errCode, errMessage})
}
