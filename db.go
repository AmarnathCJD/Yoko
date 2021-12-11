package main

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
         "fmt"
)

var (
	database = db.Database("go")
	locks_db = database.Collection("locks_db")
)

func isTrue(a string, list bson.A) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func lock_item(chat_id int64, item string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		lock := bson.D{{"chat_id", chat_id}, {"locks", []string{item}}}
		locks_db.InsertOne(context.TODO(), lock)
	} else {
                var lock_list bson.M
		locked.Decode(&lock_list)
                new_lock := lock_list["locks"].(bson.A)
                new_lock = append(new_lock, item)
                locks_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"locks", new_lock}}}})
	}
	return true
}

func get_locks(chat_id int64) bson.A {
 filter := bson.M{"chat_id": chat_id}
 locked := locks_db.FindOne(context.TODO(), filter)
 if locked.Err() != nil{
    return bson.A{}
 }
 var lock_list bson.M
 locked.Decode(&lock_list)
 fmt.Println(lock_list["locks"].(bson.A))
 return lock_list["locks"].(bson.A)
}
