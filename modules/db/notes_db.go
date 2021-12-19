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
)

func deduplicate_note(s bson.A, x string) bson.A {
	for i, v := range s {
		if v.(bson.M)["name"].(string) == x {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
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
		notez = deduplicate_note(notez, name)
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
