package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var locks_db = database.Collection("locks_dbx")

func IsTrue(a string, list bson.A) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func Lock_item(chat_id int64, items []string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		lock := bson.D{{Key: "chat_id", Value: chat_id}, {Key: "locks", Value: items}}
		locks_db.InsertOne(context.TODO(), lock)
	} else {
		var lock_list bson.M
		locked.Decode(&lock_list)
		new_lock := lock_list["locks"].(bson.A)
		for _, lock := range items {
			new_lock = append(new_lock, lock)
		}
		locks_db.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "locks", Value: new_lock}}}}, opts)
	}
	return true
}

func Unlock_item(chat_id int64, items []string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		return false
	} else {
		var lock_list bson.M
		locked.Decode(&lock_list)
		new_lock := lock_list["locks"].(bson.A)
		for _, lock := range items {
			new_lock = Remove(new_lock, lock).(bson.A)
		}
		locks_db.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "locks", Value: new_lock}}}}, opts)
	}
	return true
}

func Get_locks(chat_id int64) bson.A {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		return bson.A{}
	}
	var lock_list bson.M
	locked.Decode(&lock_list)
	return lock_list["locks"].(bson.A)
}
