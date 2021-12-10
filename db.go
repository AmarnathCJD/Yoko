package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	database = db.Database("go")
	locks_db = database.Collection("locks")
)

func lock_item(chat_id int64, item string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		fmt.Printf("locked.Err(): %v\n", locked.Err())
		lock := bson.D{{"chat_id", chat_id}, {"locks", "meaw"}}
		xd, err := locks_db.InsertOne(context.TODO(), lock)
		fmt.Println(xd)
		fmt.Printf("err.Error(): %v\n", err.Error())
	} else {
		fmt.Printf("locked: %v\n", locked)
	}
	return true
}
