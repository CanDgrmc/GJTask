package api

import (
	"github.com/CanDgrmc/gotask/models"
	"net/http"
)

type CreateRequest struct {
	*models.Maze
}

func (d *CreateRequest) Bind(r *http.Request) error {
	return nil
}
