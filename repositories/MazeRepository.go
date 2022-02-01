package repositories

import (
	"context"
	"github.com/CanDgrmc/gotask/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

const (
	collectionName string = "maze"

	mongoQueryTimeout = 10 * time.Second
)

type Repository interface {
	FindAll() ([]*models.Maze, error)
	Find(maze string) (*models.Maze, error)
	Add(maze *models.Maze) error
	Remove(id string) error
}

type MongoRepository struct {
	collection *mongo.Collection
}

func NewMazeRepository(db *mongo.Database) (*MongoRepository, error) {
	return &MongoRepository{collection: db.Collection(collectionName)}, nil
}

func (r *MongoRepository) Find(id string) (*models.Maze, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()
	var maze *models.Maze

	if err = r.collection.FindOne(ctx, bson.M{"_id": bson.M{"$eq": objectId}}).Decode(&maze); err != nil {
		return nil, err
	}

	return maze, nil
}

func (r *MongoRepository) FindAll() (*[]models.Maze, error) {
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()
	mazes := []models.Maze{}
	cursor, err := r.collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &mazes); err != nil {
		return nil, err
	}
	return &mazes, nil

}

func (r *MongoRepository) Add(m models.Maze) (*interface{}, error) {
	var (
		err    error
		result *mongo.InsertOneResult
	)
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	result, err = r.collection.InsertOne(ctx, bson.D{
		{Key: "arr", Value: m.Arr},
	})
	if err != nil {
		return nil, err
	}
	return &result.InsertedID, nil
}

func (r *MongoRepository) Delete(id string) (*int64, error) {
	var (
		err    error
		result *mongo.DeleteResult
	)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Invalid id")
	}
	ctx, cancel := context.WithTimeout(context.Background(), mongoQueryTimeout)
	defer cancel()

	if result, err = r.collection.DeleteOne(ctx, bson.M{"_id": bson.M{"$eq": objectId}}); err != nil {
		return nil, err
	}
	return &result.DeletedCount, nil
}
