package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var rules = database.Collection("rules")

func SetRules(chat_id int64, rule string) {
	rules.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "rules", Value: rules}}}}, opts)
}

func GetRules(chat_id int64) (string, string) {
	w := rules.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if w.Err() != nil {
		return "", "Rules"
	} else {
		var f bson.M
		var Rules string
		var Button = "Rules"
		w.Decode(&f)
		if co, ok := f["rules"]; ok {
			Rules = co.(string)
		}
		if bt, ok := f["btn"]; ok {
			Button = bt.(string)
		}
		return Rules, Button
	}
}

func SetRulesButton(chat_id int64, btn string) {
	rules.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "btn", Value: rules}}}}, opts)
}

func ResetRulesButton(chat_id int64) {
	rules.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "btn", Value: "Rules"}}}}, opts)
}

func PrivateRules(chat_id int64) bool {
	w := rules.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if w.Err() != nil {
		return false
	} else {
		var f bson.M
		var Private = false
		w.Decode(&f)
		if co, ok := f["private"]; ok {
			Private = co.(bool)
		}
		return Private
	}
}

func SetPrivateRules(chat_id int64, private bool) {
	rules.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "private", Value: private}}}}, opts)
}

func ResetRules(chat_id int64) {
	rules.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "rules", Value: ""}}}}, opts)
}
