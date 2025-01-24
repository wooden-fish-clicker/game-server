package db

import (
	"context"

	"game-server/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// var MongoDB *mongo.Client

type MongoDB struct {
	Client *mongo.Client
}

func NewMongoDBClient(connString string) *MongoDB {
	mongoClient := connectMongoDB(connString)
	return &MongoDB{mongoClient}
}

func connectMongoDB(connString string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(connString)

	mongoClient, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		logger.Fatal("無法連接到MongoDB: ", err)
	}

	err = mongoClient.Ping(context.TODO(), nil)
	if err != nil {
		logger.Fatal("無法連接到MongoDB: ", err)
	}

	return mongoClient

}

func (m *MongoDB) CloseMongoDB() {
	defer m.Client.Disconnect(context.TODO())
}
