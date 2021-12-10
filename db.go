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
                var lock_list mapType
		locked.Decode(&lock_list)
                new_locks := lock_list["locks"].([]interface{})
                new_locks = append(new_locks, item)
                _, err := locks_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"locks", new_locks}}}})
                fmt.Println(err)
	}
	return true
}

func get_locks(chat_id int64) string{
 filter := bson.M{"chat_id": chat_id}
 locked := locks_db.FindOne(context.TODO(), filter)
 if locked.Err() != nil{
    return ""
 }
 var lock_list bson.M
 locked.Decode(&lock_list)
 return fmt.Sprint(lock_list["locks"].([]interface{}))
}
