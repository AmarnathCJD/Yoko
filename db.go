package main

import (
	"context"
        "go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
        "github.com/google/uuid"
)

var (
	database = db.Database("go")
        opts = options.Update().SetUpsert(true)
	locks_db = database.Collection("locks_dbx")
	notes_db = database.Collection("notde")
        feds = database.Collection("feda")
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

func deduplicate_note(s bson.A, x string) bson.A {
	for i, v := range s {
		if v.(bson.M)["name"].(string) == x {
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
		locks_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"locks", new_lock}}}}, opts)
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
		locks_db.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"locks", new_lock}}}}, opts)
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

func make_fed(user_id int64, fedname string) {
 uid := uuid.New().String()
 fmt.Println(uid)
}

func get_fed(user_id int64) (string, string) {
 filter := bson.M{"user_id": user_id}
 fed := feds.FindOne(context.TODO(), filter)
 if fed.Err() != nil {
    return false, "", ""
 }
 var fed_info bson.M
 fed.Decode(&fed_info)
 return true, fed_info["fed_id"].(string), fed_info["fedname"].(string)
}
