package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	warns    = database.Collection("warns")
	settings = database.Collection("warn_settings")
)

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

func Warn_user(chat_id int64, user_id int64, reason string) (bool, int32, int32) {
	filter := bson.M{"chat_id": chat_id, "user_id": user_id}
	w := warns.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		warner := bson.M{"chat_id": chat_id, "user_id": user_id, "reasons": []string{reason}, "count": 1}
		warns.InsertOne(context.TODO(), warner)
		return false, 3, 1
	} else {
		var ww bson.M
		w.Decode(&ww)
		reasons := ww["reasons"].(bson.A)
		reasons = append(reasons, reason)
		count := ww["count"].(int32)
		count++
		warns.UpdateOne(context.TODO(), filter, bson.D{{Key: "set", Value: bson.D{{Key: "reasons", Value: reasons}, {Key: "count", Value: count}}}}, opts)
		limit := int32(3)
		if a, i := IndexInSlice(WARN_SETTINGS, "chat_id", chat_id); a {
			limit = WARN_SETTINGS[i].(bson.M)["limit"].(int32)
		}
		if count >= limit {
			warns.UpdateOne(context.TODO(), filter, bson.D{{Key: "set", Value: bson.D{{Key: "count", Value: 0}}}}, opts)
			return true, limit, count
		} else {
			return false, limit, count
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
		warns.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "reasons", Value: reasons}, {Key: "count", Value: count}}}}, opts)
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
	w := settings.FindOne(context.TODO(), filter)
	mode, time := "ban", int32(0)
	if w.Err() == nil {
		var warn bson.M
		w.Decode(&warn)
		if t, ok := warn["time"]; ok && t != nil {
			time = t.(int32)
		}
		if m, ok := warn["mode"]; ok && m != nil {
			mode = m.(string)
		}
	}
	settings.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "limit", Value: limit}, {Key: "mode", Value: mode}, {Key: "time", Value: time}}}}, opts)
}

func Set_warn_mode(chat_id int64, mode string, time int) {
	filter := bson.M{"chat_id": chat_id}
	w := settings.FindOne(context.TODO(), filter)
	limit := int32(3)
	if w.Err() == nil {
		var warn bson.M
		w.Decode(&warn)
		if l, ok := warn["limit"]; ok && l != nil {
			limit = l.(int32)
		}
	}
	settings.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "mode", Value: mode}, {Key: "time", Value: int32(time)}, {Key: "limit", Value: limit}}}}, opts)
}

func Get_warn_settings(chat_id int64) (int32, string, int32) {
	mode, limit, time := "ban", int32(3), int32(0)
	filter := bson.M{"chat_id": chat_id}
	w := settings.FindOne(context.TODO(), filter)
	if w.Err() != nil {
		return limit, mode, time
	} else {
		var warn bson.M
		w.Decode(&warn)
		if lt, ok := warn["limit"]; ok {
			limit = lt.(int32)
		}
		if md, ok := warn["mode"]; ok {
			mode = md.(string)
		}
		if tm, ok := warn["time"]; ok {
			time = tm.(int32)
		}
		return limit, mode, time
	}
}
