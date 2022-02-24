package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var cb = database.Collection("chatbot")

var CB_CHATS = _load_chats()

func _load_chats() []int64 {
	var files []bson.M
	r, _ := cb.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &files)
	array := []int64{}
	for _, x := range files {
		array = append(array, x["chat_id"].(int64))
	}
	return array
}

func SetCHatBotMode(chat_id int64, mode bool) {
	cb.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.M{"$set": bson.M{"mode": mode}}, opts)
	if mode {
		CB_CHATS = append(CB_CHATS, chat_id)
	} else {
		CB_CHATS = Remove(CB_CHATS, chat_id).([]int64)
	}
}

func GetChatbotMode(chat_id int64) bool {
	c := cb.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if c.Err() != nil {
		return false
	} else {
		var s bson.M
		c.Decode(&s)
		return s["mode"].(bool)
	}
}

func IsChatbot(chat_id int64) bool {
	for _, x := range CB_CHATS {
		if x == chat_id {
			return true
		}
	}
	return false
}
