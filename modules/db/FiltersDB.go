package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var filters = database.Collection("filters")

var Filters = make(map[int64][]MsgDB)

func SaveFilter(chatID int64, msg MsgDB) error {
	var FilterDB bson.M
	var Filter []MsgDB
	if F := filters.FindOne(context.TODO(), bson.M{"chat_id": chatID}); F.Err() != nil {
		F.Decode(&Filters)
		if F, ok := FilterDB["filters"]; ok {
			Filter = F.([]MsgDB)
		}
	}
	Filter = DupFunc(Filter, msg.Name)
	Filter = append(Filter, msg)
	_, err := filters.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$set", Value: bson.D{{Key: "filters", Value: Filters}}}}, opts)
	Filters[chatID] = Filter
	return err
}

func GetFilters(chatID int64) []MsgDB {
	if F, ok := Filters[chatID]; ok {
		return F
	} else {
		return nil
	}
}

func GetFilter(chat_id int64, name string) *MsgDB {
	if F, ok := Filters[chat_id]; ok {
		for _, x := range F {
			if x.Name == name {
				return &x
			}
		}
	}
	return nil
}

func DelFilter(chatID int64, Name string) error {
	var Fl []MsgDB
	if F, ok := Filters[chatID]; ok {
		Filters[chatID] = DupFunc(F, Name)
		Fl = Filters[chatID]
	}
	_, err := filters.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$set", Value: bson.D{{Key: "filters", Value: Fl}}}}, opts)
	return err
}

func DelAllFilters(chatID int64) error {
	_, err := filters.DeleteOne(context.TODO(), bson.M{"chat_id": chatID})
	Filters[chatID] = []MsgDB{}
	return err
}

func FilterExists(chatID int64, name string) bool {
	if F, ok := Filters[chatID]; ok {
		for _, x := range F {
			if x.Name == name {
				return true
			}
		}
	}
	return false
}

func LoadFilters() map[int64][]MsgDB {
	var files []bson.M
	r, _ := filters.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &files)
	array := map[int64][]MsgDB{}
	for _, x := range files {
		for _, y := range x["filters"].(bson.A) {
			array[x["chat_id"].(int64)] = append(array[x["chat_id"].(int64)], MsgDB{y.(bson.M)["name"].(string), y.(bson.M)["text"].(string), y.(bson.M)["file"].(FileDB)})
		}
	}
	return array
}
