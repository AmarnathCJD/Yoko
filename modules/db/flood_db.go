package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var FLOOD = __load_flood_settings()

var flood = database.Collection("flood")

func SetFloodCount(chat_id int64, count int32) {
	flood.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "count", Value: count}}}}, opts)
	FLOOD[chat_id] = GetDBFlood(chat_id)
}

func SetFloodMode(chat_id int64, mode string, time int32) {
	flood.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "mode", Value: mode}, {Key: "time", Value: time}}}}, opts)
	FLOOD[chat_id] = GetDBFlood(chat_id)
}

func GetFlood(chat_id int64) SET {
	f, ok := FLOOD[chat_id]
	if !ok {
		return SET{COUNT: 0, MODE: "ban", TIME: 0}
	} else {
		return f
	}
}

func GetDBFlood(chat_id int64) SET {
	w := flood.FindOne(context.TODO(), bson.M{"chat_id": chat_id})
	if w.Err() != nil {
		return SET{COUNT: 0, MODE: "ban", TIME: 0}
	} else {
		var f bson.M
		c, m, t := int32(0), "ban", int32(0)
		w.Decode(&f)
		if co, ok := f["count"]; ok {
			c = co.(int32)
		}
		if co, ok := f["mode"]; ok {
			m = co.(string)
		}
		if co, ok := f["time"]; ok {
			t = co.(int32)
		}
		return SET{c, m, t}
	}
}

type SET struct {
	COUNT int32
	MODE  string
	TIME  int32
}

func __load_flood_settings() map[int64]SET {
	f := map[int64]SET{}
	var files []bson.M
	r, _ := flood.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &files)
	for _, x := range files {
		count, mode, time := int32(0), "ban", int32(0)
		if c, ok := x["count"]; ok {
			count = c.(int32)
		}
		if m, ok := x["mode"]; ok {
			mode = m.(string)
		}
		if t, ok := x["time"]; ok {
			time = t.(int32)
		}
		f[x["chat_id"].(int64)] = SET{COUNT: int32(count), TIME: int32(time), MODE: mode}
	}
	return f
}
