package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
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
	case bson.A:
		switch s.(type) {
		case int64:
			for i, v := range array.(bson.A) {
				array := array.(bson.A)
				if v == s.(int64) {
					return append(array[:i], array[i+1:]...)
				}
			}
		case string:
			for i, v := range array.(bson.A) {
				array := array.(bson.A)
				if v == s.(string) {
					return append(array[:i], array[i+1:]...)
				}
			}
		}
	default:
		return array
	}
	return array
}

func IndexInSlice(list bson.A, index string, value interface{}) (bool, int) {
	switch value.(type) {
	case int64:
		for i, v := range list {
			if v.(bson.M)[index].(int64) == value.(int64) {
				return true, i
			}
		}
	case string:
		for i, v := range list {
			if v.(bson.M)[index].(string) == value.(string) {
				return true, i
			}
		}
	default:
		return false, 0
	}
	return false, 0
}
