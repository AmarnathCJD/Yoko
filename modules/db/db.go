package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinit() *mongo.Client {
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://go:amar0245@lon5-c12-1.mongo.objectrocket.com:43391,lon5-c12-2.mongo.objectrocket.com:43391,lon5-c12-0.mongo.objectrocket.com:43391/go?replicaSet=24e1adf7f54a48fba7350c36009da162&retryWrites=false"))
	if err != nil {
		panic(err)
	}
	return db
}

var db = DBinit()
