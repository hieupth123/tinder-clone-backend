package database

import (
	"context"
	"fmt"
	"github.com/phamtrunghieu/tinder-clone-backend/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

var db *mongo.Database

func Init() (*mongo.Database, error) {
	if db == nil {
		cfg := config.GetConfig()
		uri := cfg.GetString("database.uri")
		database := cfg.GetString("database.db_name")

		ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
		defer cancel()
		options := options.Client()
		options.ApplyURI(uri)
		client, err := mongo.Connect(ctx, options)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		db = client.Database(database)
	}

	return db, nil
}

func GetInstance() *mongo.Database {
	return db
}
