package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	warns    = database.Collection("warns")
	settings = database.Collection("warn_settings")
)

type Warn struct {
	Chat_id int64    `bson:"chat_id"`
	User_id int64    `bson:"user_id"`
	Count   int32    `bson:"count"`
	Reasons []string `bson:"reasons"`
}

type Settings struct {
	Chat_id int64  `bson:"chat_id"`
	Limit   int32  `bson:"limit"`
	Mode    string `bson:"mode"`
	Time    int32  `bson:"time"`
}

func LoadWarnSettings() []Settings {
	var files []Settings
	r, _ := settings.Find(context.TODO(), Settings{})
	r.All(context.TODO(), &files)
	var settings []Settings
	for _, file := range files {
		settings = append(settings, Settings{Chat_id: file.Chat_id, Limit: file.Limit, Mode: file.Mode, Time: file.Time})
	}
	return settings
}

var WarnSettings = LoadWarnSettings()

func WarnUser(chat_id int64, user_id int64, reason string) (bool, int32, int32) {
	Prev := GetWarns(chat_id, user_id)
	if Prev.Chat_id == 0 {
		Prev = Warn{
			Chat_id: chat_id,
			User_id: user_id,
			Count:   0,
			Reasons: []string{},
		}
	}
	Prev.Count++
	Prev.Reasons = append(Prev.Reasons, reason)
	Settings := GetSettings(chat_id)
	if Prev.Count >= Settings.Limit {
		Prev.Count = 0
		warns.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id, "user_id": user_id}, bson.D{{Key: "$set", Value: Prev}}, opts)
		return true, Settings.Limit, Settings.Time
	} else {
		warns.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id, "user_id": user_id}, bson.D{{Key: "$set", Value: Prev}}, opts)
		return false, Settings.Limit, Settings.Time
	}

}

func GetSettings(chat_id int64) Settings {
	for _, s := range WarnSettings {
		if s.Chat_id == chat_id {
			return s
		}
	}
	return Settings{Chat_id: chat_id, Limit: 3, Mode: "ban", Time: 0}
}

func RemoveWarn(chat_id int64, user_id int64) bool {
	Prev := GetWarns(chat_id, user_id)
	if Prev.Chat_id == 0 || Prev.Count == 0 {
		return false
	} else {
		Prev.Count--
		Prev.Reasons = Prev.Reasons[1:]
		warns.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id, "user_id": user_id}, bson.D{{Key: "$set", Value: Prev}}, opts)
		return true
	}
}

func ResetWarns(chat_id int64, user_id int64) bool {
	Prev := GetWarns(chat_id, user_id)
	if Prev.Chat_id == 0 || Prev.Count == 0 {
		return false
	} else {
		Prev.Count = 0
		Prev.Reasons = []string{}
		warns.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id, "user_id": user_id}, bson.D{{Key: "$set", Value: Prev}}, opts)
		return true
	}
}

func GetWarns(chat_id int64, user_id int64) Warn {
	var warn Warn
	if r := warns.FindOne(context.TODO(), bson.M{"chat_id": chat_id, "user_id": user_id}); r.Err() == nil {
		r.Decode(&warn)
	}
	return warn
}

func ResetChatWarns(chat_id int64) {
	warns.DeleteMany(context.TODO(), Warn{Chat_id: chat_id})
}

func SetWarnLimit(chat_id int64, limit int) {
	Set := GetSettings(chat_id)
	if Set.Chat_id == 0 {
		Set = Settings{Chat_id: chat_id, Limit: int32(limit), Mode: "ban", Time: 0}
	}
	Set.Limit = int32(limit)
	settings.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: Set}}, opts)
	WarnSettings[chat_id] = Set
}

func SetWarnMode(chat_id int64, mode string, time int) {
	Set := GetSettings(chat_id)
	if Set.Chat_id == 0 {
		Set = Settings{Chat_id: chat_id, Limit: int32(3), Mode: "ban", Time: 0}
	}
	Set.Mode = mode
	Set.Time = int32(time)
	settings.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: Set}}, opts)
	WarnSettings[chat_id] = Set
}

func GetWarnSettings(chat_id int64) Settings {
	var Setting Settings
	if r := settings.FindOne(context.TODO(), bson.M{"chat_id": chat_id}); r.Err() == nil {
		r.Decode(&Setting)
	} else {
		Setting = Settings{Chat_id: chat_id, Limit: 3, Mode: "ban", Time: 0}
	}
	return Setting
}
