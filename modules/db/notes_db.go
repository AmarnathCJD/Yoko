package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	database = db.Database("go")
	opts     = options.Update().SetUpsert(true)
	notes_db = database.Collection("notde")
	pnote    = database.Collection("pnote")
)

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func deduplicate_note(s bson.A, x string) (bson.A, bool) {
	for i, v := range s {
		if v.(bson.M)["name"].(string) == x {
			return append(s[:i], s[i+1:]...), true
		}
	}
	return s, false
}

func Save_note(chat_id int64, name string, note string, file []string) bool {
	filter := bson.M{"chat_id": chat_id}
	notes := notes_db.FindOne(context.TODO(), filter)
	if notes.Err() != nil {
		var notez bson.A
		note_s := bson.M{"name": name, "note": note, "file": file}
		notez = append(notez, note_s)
		to_insert := bson.D{{"chat_id", chat_id}, {"notes", notez}}
		notes_db.InsertOne(context.TODO(), to_insert)
	} else {
		var dec_note bson.M
		notes.Decode(&dec_note)
		notez := dec_note["notes"].(bson.A)
		new_note := bson.M{"name": name, "note": note, "file": file}
		notez, _ = deduplicate_note(notez, name)
		notez = append(notez, new_note)
		notes_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"notes", notez}}}}, opts)
	}
	return true
}

func Get_notes(chat_id int64) bson.A {
	filter := bson.M{"chat_id": chat_id}
	note_find := notes_db.FindOne(context.TODO(), filter)
	var note bson.M
	note_find.Decode(&note)
	if note == nil {
		return nil
	}
	notes := note["notes"].(bson.A)
	return notes
}

func Get_note(chat_id int64, name string) bson.M {
	filter := bson.M{"chat_id": chat_id}
	note_find := notes_db.FindOne(context.TODO(), filter)
	var note bson.M
	note_find.Decode(&note)
	if note == nil {
		return nil
	}
	notes := note["notes"].(bson.A)
	for _, y := range notes {
		if y.(bson.M)["name"].(string) == name {
			return y.(bson.M)
		}
	}
	return nil
}

func Del_note(chat_id int64, name string) bool {
	filter := bson.M{"chat_id": chat_id}
	f := notes_db.FindOne(context.TODO(), filter)
	if f.Err() != nil {
		return false
	}
	var notes bson.M
	f.Decode(&notes)
	all_notes := notes["notes"].(bson.A)
	FL, rm := deduplicate_note(all_notes, name)
	notes_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"notes", FL}}}}, opts)
	return rm
}

func Set_pnote(chat_id int64, mode bool) {
	filter := bson.M{"chat_id": chat_id}
	pnote.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"mode", mode}}}}, opts)
}

func Pnote_settings(chat_id int64) bool {
	filter := bson.M{"chat_id": chat_id}
	f := pnote.FindOne(context.TODO(), filter)
	if f.Err() != nil {
		return false
	}
	var pn bson.M
	f.Decode(&pn)
	return pn["mode"].(bool)
}

func Del_all_notes(chat_id int64) {
	filter := bson.M{"chat_id": chat_id}
	notes_db.DeleteOne(context.TODO(), filter)
}
