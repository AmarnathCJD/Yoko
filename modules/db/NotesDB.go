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

type NotesDB struct {
	ChatID int64   `json:"chat_id,omitempty"`
	Notes  []MsgDB `json:"notes,omitempty"`
}

func SaveNote(chatID int64, Msg MsgDB) {
	var NoteDB NotesDB
	var Note []MsgDB
	if N := notes.FindOne(context.TODO(), bson.M{"chat_id": chatID}); N.Err() == nil {
		N.Decode(&NoteDB)
		Note = NoteDB.Notes
	}
	Note = DupFunc(Note, Msg.Name)
	Note = append(Note, Msg)
	Notes[chatID] = append(Notes[chatID], Msg)
	notes.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$set", Value: bson.D{{Key: "notes", Value: Note}}}}, opts)
}

func GetNotes(chat_id int64) []MsgDB {
	if N, ok := Notes[chat_id]; ok {
		return N
	}
	return nil
}

func GetNote(chat_id int64, name string) MsgDB {
	var NoteDB NotesDB
	var Note []MsgDB
	if N := notes.FindOne(context.TODO(), bson.M{"chat_id": chat_id}); N.Err() == nil {
		N.Decode(&NoteDB)
		Note = NoteDB.Notes
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
	Notes[chat_id] = DupFunc(NT, name)
	notes.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "notes", Value: NT}}}})
	return Exists
}

func SetPnote(chat_id int64, mode bool) {
	pnote.UpdateOne(context.TODO(), bson.M{"chat_id": chat_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "mode", Value: mode}}}}, opts)
}

func PnoteSettings(chat_id int64) bool {
	var pnotes bson.M
	if N := pnote.FindOne(context.TODO(), bson.M{"chat_id": chat_id}); N.Err() == nil {
		N.Decode(&pnotes)
		if NT, ok := pnotes["mode"]; ok {
			return NT.(bool)
		}
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
	r, _ := notes.Find(context.TODO(), bson.M{})
	r.All(context.TODO(), &files)
	array := map[int64][]MsgDB{}
	for _, x := range files {
		for _, y := range x["notes"].(bson.A) {
			y := y.(bson.M)
			var File FileDB
			var Text string
			var Name string
			if y["file"] != nil {
				FileM := y["file"].(bson.M)
				if FileM["file_id"] != nil {
					File.FileID = FileM["file_id"].(string)
				}
				if FileM["file_type"] != nil {
					File.FileType = FileM["file_type"].(string)
				}
			}
			if y["text"] != nil {
				Text = y["text"].(string)
			}
			if y["name"] != nil {
				Name = y["name"].(string)
			}
			array[x["chat_id"].(int64)] = append(array[x["chat_id"].(int64)], MsgDB{File: File, Text: Text, Name: Name})
		}
	}
	return array
}
