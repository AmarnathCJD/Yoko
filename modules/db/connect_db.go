package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var c = database.Collection("connect")

func ConnectChat(chatID int64, userID int64) {
	c.UpdateOne(context.TODO(), bson.M{"user_id": userID}, bson.M{"$set": bson.M{"chat_id": chatID}})
}

func GetChat(userID int64) int64 {
	var chatID int64
	chat := c.FindOne(context.TODO(), bson.M{"user_id": userID})
	if chat.Err() != nil {
		return 0
	}
	chat.Decode(&chatID)
	return chatID
}
