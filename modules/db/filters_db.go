package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var filters = database.Collection("filters")

var FILTERS = _load_filters()

func deduplicate_filters(s []string, x string) ([]string, bool, int) {
	for i, v := range s {
		if v == x {
			return append(s[:i], s[i+1:]...), true, i
		}
	}
	return s, false, 0
}

func _load_filters() map[int64][]string {
	var files []bson.M
	r, _ := filters.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &files)
	array := map[int64][]string{}
	for _, x := range files {
		for _, y := range x["filters"].(bson.A) {
			array[x["chat_id"].(int64)] = append(array[x["chat_id"].(int64)], y.(bson.M)["name"].(string))
		}
	}
	return array
}

func Save_filter(chat_id int64, name string, note string, file []string) bool {
	filter := bson.M{"chat_id": chat_id}
	fill := filters.FindOne(context.TODO(), filter)
	if fill.Err() != nil {
		var f bson.A
		note_s := bson.M{"name": name, "text": note, "file": file}
		f = append(f, note_s)
		to_insert := bson.D{{Key: "chat_id", Value: chat_id}, {Key: "filters", Value: f}}
		filters.InsertOne(context.TODO(), to_insert)
	} else {
		var dec_note bson.M
		fill.Decode(&dec_note)
		f := dec_note["filters"].(bson.A)
		new_filter := bson.M{"name": name, "text": note, "file": file}
		f, _ = deduplicate_note(f, name)
		f = append(f, new_filter)
		filters.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "filters", Value: f}}}}, opts)
	}
	FILTERS[chat_id] = append(FILTERS[chat_id], name)

	return true
}

func Get_filters(chat_id int64) []string {
	a, b := FILTERS[chat_id]
	if b {
		return a
	}
	return nil
}

func Get_filter(chat_id int64, name string) bson.M {
	filter := bson.M{"chat_id": chat_id}
	f := filters.FindOne(context.TODO(), filter)
	var fil bson.M
	f.Decode(&fil)
	if fil == nil {
		return nil
	}
	fill := fil["filters"].(bson.A)
	for _, y := range fill {
		if y.(bson.M)["name"].(string) == name {
			return y.(bson.M)
		}
	}
	return nil
}

func Del_filter(chat_id int64, name string) bool {
	filter := bson.M{"chat_id": chat_id}
	fl := filters.FindOne(context.TODO(), filter)
	if fl.Err() != nil {
		return false
	}
	var f bson.M
	fl.Decode(&f)
	all_f := f["filters"].(bson.A)
	FL, rm := deduplicate_note(all_f, name)
	filters.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "filters", Value: FL}}}}, opts)
	if rm {
		FILTERS[chat_id], _, _ = deduplicate_filters(FILTERS[chat_id], name)
	}
	return rm
}

func Del_all_filters(chat_id int64) {
	filter := bson.M{"chat_id": chat_id}
	filters.DeleteOne(context.TODO(), filter)
}
