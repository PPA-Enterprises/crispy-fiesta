package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConnection struct {
	client *mongo.Client
}

func NewConnection(host string) (conn *DBConnection) {

	//TODO: Auth
	client, err := mongo.NewClient(options.Client().ApplyURI(host))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		panic(err)
	}

	conn = &DBConnection{client}
	return conn

}

func (conn *DBConnection) Use(dbName, tableName string) *mongo.Collection {
	return conn.client.Database(dbName).Collection(tableName)
}

