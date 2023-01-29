package repository

import (
	"context"
	"ecom-users/internal/config"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Models struct {
	Users UsersModel
}

func NewModels(db *mongo.Database) Models {
	return Models{
		Users: UsersModel{db: db.Collection("users")},
	}
}

func OpenDB(cfg *config.Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// "mongodb://localhost:27017/"
	URI := fmt.Sprintf("%s/%s", cfg.MONGO_URI, cfg.DB_NAME)

	clientOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil { 
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(cfg.DB_NAME), nil
}