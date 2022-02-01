package models

import (
	"github.com/CanDgrmc/gotask/lib"
	"net/http"
)

type Maze struct {
	Id  int     `json:"id,omitempty"`
	Arr [][]int `json:"arr,omitempty" bson:"arr,omitempty"`
}

func (arr *Maze) Bind(r *http.Request) error {
	return nil
}

func (arr *Maze) getLine(line int) []int {
	return arr.Arr[line]
}

func (m *Maze) ExistsInLine(needle int, line int) bool {
	return lib.ExistsIntValue(m.getLine(line), needle)
}

func (m *Maze) FindPlayerPosition(player int) *Position {

	for x := range m.Arr {
		y := lib.GetIndexOf(m.getLine(x), player)
		if y > -1 {
			return &Position{X: x, Y: y}
		}
	}
	return nil
}

func (m *Maze) isRectangular() bool {
	return len(m.Arr) != len(m.getLine(0))
}

func (m *Maze) Validate(maks int) bool {
	if !m.isRectangular() {
		return false
	}

	totalCount := len(m.Arr)
	for i := range m.Arr {
		totalCount += len(m.getLine(i))
	}

	if totalCount > maks {
		return false
	}

	return true
}
