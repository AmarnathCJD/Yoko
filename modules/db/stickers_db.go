package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var stickers = database.Collection("stick")

func Add_sticker(user_id int64, name string) {
	filter := bson.M{"user_id": user_id}
	s := stickers.FindOne(context.TODO(), filter)
	if s.Err() != nil {
		var packs bson.A
		packs = append(packs, bson.M{"name": name, "count": 1})
		stickers.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"packs", packs}}}}, opts)
	} else {
		var stick bson.M
		s.Decode(&stick)
		packs := stick["packs"].(bson.A)
		packs = append(packs, bson.M{"name": name, "count": 1})
		stickers.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"packs", packs}}}}, opts)
	}

}

func Get_user_pack(user_id int64) (bool, int32, string) {
	filter := bson.M{"user_id": user_id}
	s := stickers.FindOne(context.TODO(), filter)
	if s.Err() != nil {
		return false, 0, ""
	} else {
		var stick bson.M
		s.Decode(&stick)
		packs := stick["packs"].(bson.A)
		return true, packs[len(packs)-1].(bson.M)["count"].(int32), packs[len(packs)-1].(bson.M)["name"].(string)
	}
}

func Update_count(user_id int64, name string) {
	filter := bson.M{"user_id": user_id}
	s := stickers.FindOne(context.TODO(), filter)
	var stick bson.M
	s.Decode(&stick)
	packs := stick["packs"].(bson.A)
	c := packs[len(packs)-1].(bson.M)["count"].(int)
	c++
	packs[len(packs)-1] = bson.M{"name": packs[len(packs)-1].(bson.M)["name"].(string), "count": c}
	stickers.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"packs", packs}}}}, opts)
}
