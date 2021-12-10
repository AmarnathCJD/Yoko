package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	database = db.Database("go")
	locks_db = database.Collection("locks_db")
)

func lock_item(chat_id int64, item string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		lock := bson.D{{"chat_id", chat_id}, {"locks", []string{item}}}
		locks_db.InsertOne(context.TODO(), lock)
	} else {
                var lock_list bson.M
		locks := locked.Decode(&lock_list)
                fmt.Println(lock_list)
	}
	return true
}
