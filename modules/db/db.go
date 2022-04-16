package db

import (
	"context"
        "os"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database = db.Database("go")
	opts     = options.Update().SetUpsert(true)
)

func DBinit() *mongo.Client {
	db, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URI"))
	if err != nil {
		panic(err)
	}
	return db
}

type MsgDB struct {
	Name string `json:"name,omitempty"`
	Text string `json:"text,omitempty"`
	File FileDB `json:"file,omitempty"`
}
type FileDB struct {
	FileID   string `json:"file_id,omitempty"`
	FileType string `json:"file_type,omitempty"`
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

func DupFunc(F []MsgDB, Name string) []MsgDB {
	for i, x := range F {
		if x.Name == Name {
			return append(F[:i], F[i+1:]...)
		}
	}
	return F
}
