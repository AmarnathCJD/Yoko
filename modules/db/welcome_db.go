package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	welcome         = database.Collection("welcum")
	default_welcome = "Hi %s, welcome to %s!"
)

func Set_welcome(chat_id int64, text string, file []string) {
	welcome.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.M{"$set": bson.M{"text": text, "file": file}}, opts)
}

func Set_welcome_mode(chat_id int64, mode bool) {
	welcome.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.M{"$set": bson.M{"mode": mode}}, opts)
}

func Get_welcome(chat_id int64) (string, bson.A, bool) {
	text, file, mode := "", bson.A{}, true
	w := welcome.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if w.Err() != nil {
		return default_welcome, file, mode
	} else {
		var welcome_settings bson.M
		w.Decode(&welcome_settings)
		if md, ok := welcome_settings["mode"]; ok {
			mode = md.(bool)
		}
		if txt, ok := welcome_settings["text"]; ok {
			text = txt.(string)
		}
		if fl, ok := welcome_settings["file"]; ok {
			file = fl.(bson.A)
		}
	}
	return text, file, mode
}

func Reset_welcome(chat_id int64) {
	welcome.DeleteOne(context.TODO(), bson.M{"chat_id": chat_id})
}
