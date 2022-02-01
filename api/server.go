package api

import (
	"context"
	"github.com/CanDgrmc/gotask/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/viper"
)

type Server struct {
	*http.Server
}

func NewServer(client *mongo.Client, config *Configuration) (*Server, error) {

	db := client.Database(config.Database)

	mazeRepository, err := repositories.NewMazeRepository(db)

	if err != nil {
		return nil, err
	}

	api, err := New(viper.GetBool("enable_cors"), config, mazeRepository)
	if err != nil {
		return nil, err
	}

	port := config.Port
	addr := ":" + port

	srv := http.Server{
		Addr:    addr,
		Handler: api,
	}

	return &Server{&srv}, nil
}

func (srv *Server) Start() {
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()
	log.Printf("Listening on %s\n", srv.Addr)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	sig := <-quit
	log.Println("Shutting down server... Reason:", sig)

	if err := srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	log.Println("stopped")
}
