package db

import (
 "context"
 "go.mongodb.org/mongo-driver/bson"
 "fmt"
)

var warns = database.Collection("warns")

func Warn_user(chat_id int64, user_id int64, reason string) {
 filter := bson.M{"chat_id": chat_id, "user_id": user_id}
 w := warns.FindOne(context.TODO(), filter)
 if w.Err() != nil {
    warner := bson.M
    fmt.Println(warner)
 }
}
 
