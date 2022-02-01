package api

import (
	"errors"
	"fmt"
	"github.com/CanDgrmc/gotask/models"
	"github.com/CanDgrmc/gotask/repositories"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func New(enableCORS bool, config *Configuration, mazeRepository *repositories.MongoRepository) (*chi.Mux, error) {

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Timeout(15 * time.Second))

	r.Use(render.SetContentType(render.ContentTypeJSON))

	if enableCORS {
		r.Use(corsConfig().Handler)
	}
	r.Get("/maze", func(w http.ResponseWriter, r *http.Request) {
		listMazes(w, r, mazeRepository)
	})
	r.Post("/maze", func(w http.ResponseWriter, r *http.Request) {
		createMaze(w, r, config, mazeRepository)
	})
	r.Delete("/maze/{id}", func(w http.ResponseWriter, r *http.Request) {
		deleteMaze(w, r, mazeRepository)
	})
	r.Get("/maze/{id}", func(w http.ResponseWriter, r *http.Request) {
		getMaze(w, r, config, mazeRepository)
	})
	r.Get("/maze/{id}/calculate", func(w http.ResponseWriter, r *http.Request) {
		calculateMaze(w, r, config, mazeRepository)
	})

	return r, nil
}

func corsConfig() *cors.Cors {

	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           86400, // Maximum value not ignored by any of major browsers
	})
}

func listMazes(w http.ResponseWriter, r *http.Request, mazeRepository *repositories.MongoRepository) {
	var (
		mazes *[]models.Maze
		err   error
	)
	mazes, err = mazeRepository.FindAll()

	if err != nil {
		render.Respond(w, r, nil)
	}
	render.Respond(w, r, mazes)

}

func createMaze(w http.ResponseWriter, r *http.Request, config *Configuration, mazeRepository *repositories.MongoRepository) {
	var (
		insertedId interface{}
		err        error
	)
	var maze models.Maze
	data := &CreateRequest{}

	if err = render.Bind(r, data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	maze.Arr = data.Arr
	if !maze.Validate(config.Maze.Maks_element_count) {
		http.Error(w, "invalid-payload", http.StatusBadRequest)
	}

	if insertedId, err = mazeRepository.Add(maze); err != nil {
		w.Write([]byte(err.Error()))
	}
	render.Respond(w, r, insertedId)
}

func getMaze(w http.ResponseWriter, r *http.Request, config *Configuration, mazeRepository *repositories.MongoRepository) {

	id := chi.URLParam(r, "id")

	maze, err := mazeRepository.Find(id)
	if err != nil {
		render.Respond(w, r, err.Error())
		return
	}
	render.Respond(w, r, maze)
}
func deleteMaze(w http.ResponseWriter, r *http.Request, mazeRepository *repositories.MongoRepository) {
	id := chi.URLParam(r, "id")

	res, err := mazeRepository.Delete(id)
	if err != nil {
		render.Respond(w, r, err.Error())
		return
	}
	render.Respond(w, r, res)
}

func calculateMaze(w http.ResponseWriter, r *http.Request, config *Configuration, mazeRepository *repositories.MongoRepository) {

	id := chi.URLParam(r, "id")
	var visited []models.Position
	var fails []models.Position
	var solutions [][]models.Position
	maze, err := mazeRepository.Find(id)

	if err != nil {
		render.Respond(w, r, err.Error())
		return
	}

	playerPosition := maze.FindPlayerPosition(config.Maze.Player)
	for {
		if solution, err := solve(playerPosition, maze, config.Maze, visited, maze.FindPlayerPosition(config.Maze.Player), fails, solutions); err != nil {
			break
		} else {
			solutions = append(solutions, *solution)
			playerPosition = maze.FindPlayerPosition(config.Maze.Player)
		}
	}
	fmt.Println(solutions)
	shortestSolution := getShortest(solutions)

	render.Respond(w, r, shortestSolution)
	return
}

func solve(playerPosition *models.Position, maze *models.Maze, mazeConfig MazeConfig, visited []models.Position, startPosition *models.Position, fails []models.Position, solutions [][]models.Position) (*[]models.Position, error) {
	if maze.Arr[playerPosition.X][playerPosition.Y] == mazeConfig.Space && playerPosition.X == 0 {
		return &visited, nil
	}
	nextTop := playerPosition.GetNextTop()
	nextRight := playerPosition.GetNextRight()
	nextLeft := playerPosition.GetNextLeft()
	if isSafe(nextTop, playerPosition.Y, maze.Arr, mazeConfig) && !isVisited(nextTop, playerPosition.Y, visited) && !isFailure(nextTop, playerPosition.Y, fails) && !repeatedSolution(solutions, nextTop, playerPosition.Y, len(visited)) {
		fmt.Println("moving top..")
		playerPosition.MoveTop()
		visited = append(visited, models.Position{X: playerPosition.X, Y: playerPosition.Y})

	} else if isSafe(playerPosition.X, nextRight, maze.Arr, mazeConfig) && !isVisited(playerPosition.X, nextRight, visited) && !isFailure(playerPosition.X, nextRight, fails) && !repeatedSolution(solutions, playerPosition.X, nextRight, len(visited)) {
		fmt.Println("moving right..")
		playerPosition.MoveRight()
		visited = append(visited, models.Position{X: playerPosition.X, Y: playerPosition.Y})

	} else if isSafe(playerPosition.X, nextLeft, maze.Arr, mazeConfig) && !isVisited(playerPosition.X, nextLeft, visited) && !isFailure(playerPosition.X, nextLeft, fails) && !repeatedSolution(solutions, playerPosition.X, nextLeft, len(visited)) {
		fmt.Println("moving left..")
		playerPosition.MoveLeft()
		visited = append(visited, models.Position{X: playerPosition.X, Y: playerPosition.Y})
	} else {
		if playerPosition.X == startPosition.X && playerPosition.Y == startPosition.Y {
			return nil, errors.New("done")
		}
		beforeFail := pop2(&visited)
		fails = append(fails, models.Position{X: playerPosition.X, Y: playerPosition.Y})

		playerPosition.X = beforeFail.X
		playerPosition.Y = beforeFail.Y

	}

	return solve(playerPosition, maze, mazeConfig, visited, startPosition, fails, solutions)

}

func isSafe(x int, y int, maze [][]int, mazeConfig MazeConfig) bool {
	return maze[x][y] == mazeConfig.Space
}

func isVisited(x int, y int, visited []models.Position) bool {
	for _, i := range visited {
		if i.X == x && i.Y == y {
			return true
		}
	}
	return false
}

func isFailure(x int, y int, failure []models.Position) bool {

	for _, i := range failure {
		if i.X == x && i.Y == y {
			fmt.Println("fail")
			return true
		}
	}
	return false
}

func repeatedSolution(solutions [][]models.Position, x int, y int, key int) bool {
	for _, sol := range solutions {
		if len(sol) > key && sol[key].X == x && sol[key].Y == y {
			fmt.Println("repeat")
			return true
		}
	}
	return false
}

func pop2(alist *[]models.Position) models.Position {
	f := len(*alist)

	rv := (*alist)[f-2]
	*alist = (*alist)[:f-2]
	return rv
}

func getShortest(solutions [][]models.Position) []models.Position {
	shortestIndex := -1
	for i, el := range solutions {
		if shortestIndex == -1 || len(solutions[shortestIndex]) > len(el) {
			shortestIndex = i
		}
	}
	return solutions[shortestIndex]
}
