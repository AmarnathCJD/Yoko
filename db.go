package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	database = db.Database("go")
	locks_db = database.Collection("locks_dbx")
	notes_db = database.Collection("notes")
)

func isTrue(a string, list bson.A) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func remove(s bson.A, r string) bson.A {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func lock_item(chat_id int64, items []string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		lock := bson.D{{"chat_id", chat_id}, {"locks", items}}
		locks_db.InsertOne(context.TODO(), lock)
	} else {
		var lock_list bson.M
		locked.Decode(&lock_list)
		new_lock := lock_list["locks"].(bson.A)
		for _, lock := range items {
			new_lock = append(new_lock, lock)
		}
		locks_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"locks", new_lock}}}})
	}
	return true
}

func unlock_item(chat_id int64, items []string) bool {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		return false
	} else {
		var lock_list bson.M
		locked.Decode(&lock_list)
		new_lock := lock_list["locks"].(bson.A)
		for _, lock := range items {
			new_lock = remove(new_lock, lock)
		}
		locks_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"locks", new_lock}}}})
	}
	return true
}

func get_locks(chat_id int64) bson.A {
	filter := bson.M{"chat_id": chat_id}
	locked := locks_db.FindOne(context.TODO(), filter)
	if locked.Err() != nil {
		return bson.A{}
	}
	var lock_list bson.M
	locked.Decode(&lock_list)
	return lock_list["locks"].(bson.A)
}

func save_note(chat_id int64, name string, note string, file []string) bool {
	filter := bson.M{"chat_id": chat_id}
	notes := notes_db.FindOne(context.TODO(), filter)
	if notes.Err() != nil {
		var notes bson.A
		note := bson.M{"name": name, "note": note, "file": file}
		fmt.Println(note)
		notes = append(notes, note)
		to_insert := bson.D{{"chat_id", chat_id}, {"notes", notes}}
		notes_db.InsertOne(context.TODO(), to_insert)
	} else {
		var dec_note bson.M
		notes.Decode(&dec_note)
		note := dec_note["notes"].(bson.A)
		new_note := bson.M{"name": name, "note": note, "file": file}
		note = append(note, new_note)
		notes_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"notes", note}}}})
	}
	return true
}

func get_notes(chat_id int64) bson.A {
	filter := bson.M{"chat_id": chat_id}
	note_find := notes_db.FindOne(context.TODO(), filter)
	var note bson.M
	note_find.Decode(&note)
        if note == nil{
           return nil
        }
	notes := note["notes"].(bson.A)
	return notes
}

func get_note(chat_id int64, name string) bson.M {
 filter := bson.M{"chat_id": chat_id}
 note_find := notes_db.FindOne(context.TODO(), filter)
 var note bson.M
 note_find.Decode(&note)
 if note == nil{
           return nil
        }
 notes := note["notes"].(bson.A)
 for _, y := range notes{
   if y.(bson.M)["name"].(string) == name{
      return y.(bson.M)
   }
 }
 return nil
}
