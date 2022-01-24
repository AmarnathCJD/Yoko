package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var cb = database.Collection("chatbot")

var CHATS = _load_chats()

func removeInt64(s []int64, r int64) []int64 {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

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

func Set_chatbot_mode(chat_id int64, mode bool) {
	cb.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.M{"$set": bson.M{"mode": mode}}, opts)
	if !mode {
		CHATS = append(CHATS, chat_id)
	} else {
		CHATS = removeInt64(CHATS, chat_id)
	}
}

func Get_chatbot_mode(chat_id int64) bool {
	c := cb.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if c.Err() != nil {
		return false
	} else {
		var s bson.M
		c.Decode(&s)
		return s["mode"].(bool)
	}
}
