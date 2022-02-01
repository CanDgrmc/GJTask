package api

import "github.com/CanDgrmc/gotask/models"

type GetAllResponse struct {
	Mazes   *[]models.Maze
	Success bool `json:"success"`
}
