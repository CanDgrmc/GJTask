package cmd

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"

	"github.com/CanDgrmc/gotask/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server ",
	Run:   run,
}

func init() {
	RootCmd.AddCommand(serveCmd)
}

func run(cmd *cobra.Command, args []string) {
	var config *api.Configuration
	viper.Unmarshal(&config)

	clientOptions := options.Client().
		ApplyURI(config.Mongo_connection_string)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mongo is connected")

	server, err := api.NewServer(client, config)

	if err != nil {
		log.Fatal(err)
	}
	server.Start()
}
