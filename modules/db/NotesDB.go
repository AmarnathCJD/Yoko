package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	notes = database.Collection("notde")
	pnote = database.Collection("pnote")
	Notes = LoadNotes()
)

func SaveNote(chatID int64, Msg MsgDB) {
	var NotesDB bson.M
	var Note []MsgDB
	if N := notes.FindOne(context.TODO(), bson.M{"chat_id": chatID}); N.Err() != nil {
		N.Decode(&Filters)
                if NT, ok := NotesDB["notes"] ; ok {
		Note = NT.([]MsgDB)
}
	}
	Note = DupFunc(Note, Msg.Name)
	Note = append(Note, Msg)
	Notes[chatID] = Note
	notes.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$set", Value: bson.D{{Key: "notes", Value: Note}}}})
}

func GetNotes(chat_id int64) []MsgDB {
	if N, ok := Notes[chat_id]; ok {
		return N
	}
	return nil
}

func GetNote(chat_id int64, name string) MsgDB {
	var NotesDB bson.M
	var Note []MsgDB
	if N := notes.FindOne(context.TODO(), bson.M{"chat_id": chat_id}); N.Err() != nil {
		N.Decode(&NotesDB)
		if NT, ok := NotesDB["notes"] ; ok {
		Note = NT.([]MsgDB)
}
		for _, y := range Note {
			if y.Name == name {
				return y
			}
		}
	}
	return MsgDB{}
}

func DelNote(chat_id int64, name string) bool {
	var NT []MsgDB
	var Exists bool
	if N, ok := Notes[chat_id]; ok {
		NT = N
	}
	for i, x := range NT {
		if x.Name == name {
			NT = append(NT[:i], NT[i+1:]...)
			Exists = true
			break
		}
	}
	Notes[chat_id] = NT
	notes.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "notes", Value: NT}}}})
	return Exists
}

func SetPnote(chat_id int64, mode bool) {
	pnote.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "mode", Value: mode}}}}, opts)
}

func PnoteSettings(chat_id int64) bool {
	var pnotes bson.M
	if N := pnote.FindOne(context.TODO(), bson.M{"chat_id": chat_id}); N.Err() != nil {
		N.Decode(&pnotes)
		return pnotes["mode"].(bool)
	}
	return false
}

func DelAllNotes(chat_id int64) {
	notes.DeleteOne(context.TODO(), bson.M{"chat_id": chat_id})
	Notes[chat_id] = []MsgDB{}
}

func NoteExists(chat_id int64, name string) bool {
	if N, ok := Notes[chat_id]; ok {
		for _, x := range N {
			if x.Name == name {
				return true
			}
		}
	}
	return false
}

func LoadNotes() map[int64][]MsgDB {
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
