package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var warns = database.Collection("warns")
var settings = database.Collection("warn_settings")

func __load_warn_settings() bson.A {
	var files []bson.M
	r, _ := settings.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &files)
	array := bson.A{}
	for _, x := range files {
		array = append(array, x)
	}
	return array
}

func IndexInSlice(list bson.A, index string, value int64) (bool, int) {
	for i, x := range list {
		if x.(bson.M)[index] == value {
			return true, i
		}
	}
	return false, 0
}

var WARN_SETTINGS = __load_warn_settings()

func Warn_user(chat_id int64, user_id int64, reason string) (bool, int32) {
	filter := bson.M{"chat_id": chat_id, "user_id": user_id}
	w := warns.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		warner := bson.M{"chat_id": chat_id, "user_id": user_id, "reasons": []string{reason}, "count": 1}
		warns.InsertOne(context.TODO(), warner)
		return false, 3
	} else {
		var ww bson.M
		w.Decode(&ww)
		reasons := ww["reasons"].(bson.A)
		reasons = append(reasons, reason)
		count := ww["count"].(int32)
		count++
		warns.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"reasons", reasons}, {"count", count}}}}, opts)
		limit := int32(3)
		if a, i := IndexInSlice(WARN_SETTINGS, "chat_id", chat_id); a {
			limit = WARN_SETTINGS[i].(bson.M)["limit"].(int32)
		}
		if count >= limit {
			warns.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"count", 0}}}}, opts)
			return true, limit
		} else {
			return false, limit
		}
	}
}

func Remove_warn(chat_id int64, user_id int64) bool {
	filter := bson.M{"chat_id": chat_id, "user_id": user_id}
	w := warns.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		return false
	} else {
		var ww bson.M
		w.Decode(&ww)
		count := ww["count"].(int32)
		if count == 0 {
			return false
		}
		count--
		reasons := ww["reasons"].(bson.A)
		reasons = reasons[:len(reasons)-1]
		warns.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"reasons", reasons}, {"count", count}}}}, opts)
		return true
	}
}

func Reset_warns(chat_id int64, user_id int64) bool {
	filter := bson.M{"chat_id": chat_id, "user_id": user_id}
	w := warns.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		return false
	} else {
		warns.DeleteOne(context.TODO(), filter)
		return true
	}
}

func Get_warns(chat_id int64, user_id int64) bson.M {
	filter := bson.M{"chat_id": chat_id, "user_id": user_id}
	w := warns.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		return nil
	} else {
		var ww bson.M
		w.Decode(&ww)
		return ww
	}
}

func Reset_chat_warns(chat_id int64) {
	filter := bson.M{"chat_id": chat_id}
	warns.DeleteMany(context.TODO(), filter)
}

func Set_warn_limit(chat_id int64, limit int) {
	filter := bson.M{"chat_id": chat_id}
	w := warns.FindOne(context.TODO(), filter)
	mode, time := "ban", int32(0)
	if w.Err() == nil {
		var warn bson.M
		w.Decode(&warn)
		mode, time = warn["mode"].(string), warn["time"].(int32)
	}
	settings.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{"$set", bson.D{{"limit", limit}, {"mode", mode}, {"time", time}}}}, opts)
}

func Set_warn_mode(chat_id int64, mode string, time int) {
	filter := bson.M{"chat_id": chat_id}
	w := warns.FindOne(context.TODO(), filter)
	limit := int32(3)
	if w.Err() == nil {
		var warn bson.M
		w.Decode(&warn)
		limit = warn["limit"].(int32)
	}
	settings.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{"$set", bson.D{{"mode", mode}, {"time", int32(time)}, {"limit", limit}}}}, opts)
}

func Get_warn_settings(chat_id int64) (int32, string, int32) {
	filter := bson.M{"chat_id": chat_id}
	w := warns.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		return 3, "ban", 0
	} else {
		var warn bson.M
		w.Decode(&warn)
		return warn["limit"].(int32), warn["mode"].(string), warn["time"].(int32)
	}
}
