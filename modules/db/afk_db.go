package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

var afk = database.Collection("afk")

func Set_afk(user_id int64, fname string, reason string) {
	afk.UpdateOne(context.TODO(), bson.M{"user_id": user_id}, bson.M{"fname": fname, "reason": reason, "time": time.Now().Unix()}, opts)
	AFK = append(AFK, user_id)
}

func Unset_afk(user_id int64) {
	afk.DeleteOne(context.TODO(), bson.M{"user_id": user_id})
	AFK = removeInt(AFK, user_id)
}

func IsAfk(user_id int64) bool {
	for _, x := range AFK {
		if x.(int64) == user_id {
			return true
		}
	}
	return false
}

func GetAfk(user_id int64) bson.M {
	u := afk.FindOne(context.TODO(), bson.M{"user_id": user_id})
	if u.Err() == nil {
		var det bson.M
		u.Decode(&det)
		return det
	}
	return nil
}

func __load_afk() bson.A {
	var afks []bson.M
	var AFK bson.A
	r, _ := afk.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &afks)
	for _, x := range afks {
		AFK = append(AFK, x["user_id"].(int64))
	}
	return AFK
}

var AFK = __load_afk()
