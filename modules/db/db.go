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

func Remove(array interface{}, s interface{}) interface{} {
	switch array.(type) {
	case []int64:
		for i, v := range array.([]int64) {
			array := array.([]int64)
			if v == s.(int64) {
				return append(array[:i], array[i+1:]...)
			}
		}
	case []string:
		for i, v := range array.([]string) {
			array := array.([]string)
			if v == s.(string) {
				return append(array[:i], array[i+1:]...)
			}
		}
	default:
		return nil
	}
	return array
}
