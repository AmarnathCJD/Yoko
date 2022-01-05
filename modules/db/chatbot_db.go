package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var cb = database.Collection("chatbot")

func Set_chatbot_mode(chat_id int64, mode bool) {
	cb.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.M{"$set": bson.M{"mode": mode}})
}

func Get_chatbot_mode(chat_id int64) bool {
	c := cb.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if c.Err() != nil {
		return false
	} else {
		var s bson.M
		c.Decode(&s)
		return c["mode"].(bool)
	}
}
