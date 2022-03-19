package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	stickers = database.Collection("stick")
)

type Pack struct {
	Name  string `json:"name"`
	Count int32  `json:"count"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Packs struct {
	UserID int64
	Packs  []Pack
}

func AddSticker(user_id int64, name string, title string, _type string) {
	Packs := GetUserPacks(user_id)
	Packs = append(Packs, Pack{name, 1, title, _type})
	stickers.UpdateOne(context.TODO(), bson.M{"user_id": user_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "packs", Value: Packs}}}}, opts)
}

func GetUserPacks(user_id int64) []Pack {
	var s Packs
	if packs := stickers.FindOne(context.TODO(), bson.M{"user_id": user_id}); packs.Err() == nil {
		packs.Decode(&s)
	}
	return s.Packs
}

func GetPack(user_id int64, _type string) (Pack, int) {
	Packs := GetUserPacks(user_id)
	var p Pack
	Q := 0
	for _, x := range Packs {
		if x.Type == _type {
			p = x
			Q++
		}
	}
	return p, Q
}

func UpdateCount(user_id int64, _type string) {
	Packs := GetUserPacks(user_id)
	var p Pack
	var Index int
	for i, x := range Packs {
		if x.Type == _type {
			p = x
			Index = i
		}
	}
	p.Count++
	Packs[Index] = p
	stickers.UpdateOne(context.TODO(), bson.M{"user_id": user_id}, bson.D{{Key: "$set", Value: bson.D{{Key: "packs", Value: Packs}}}}, opts)
}
