package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var stickers = database.Collection("stick")

func Add_sticker(user_id int64, name string, title string, _type string) {
	filter := bson.M{"user_id": user_id, "type": _type}
	s := stickers.FindOne(context.TODO(), filter)
	if s.Err() != nil {
		var packs bson.A
		packs = append(packs, bson.M{"name": name, "count": 1, "title": title})
		stickers.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "packs", Value: packs}}}}, opts)
	} else {
		var stick bson.M
		s.Decode(&stick)
		packs := stick["packs"].(bson.A)
		packs = append(packs, bson.M{"name": name, "count": 1, "title": title})
		stickers.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "packs", Value: packs}}}}, opts)
	}

}

func Get_user_pack(user_id int64, _type string) (bool, int32, string) {
	filter := bson.M{"user_id": user_id, "type": _type}
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

func Get_user_packs(user_id int64) map[string]string {
	var files []bson.M
	filter := bson.M{"user_id": user_id}
	r, _ := filters.Find(context.TODO(), filter)
	r.All(context.TODO(), &files)
	fmt.Println(files)
	if len(files) == 0 {
		return nil
	} else {
		var Names = make(map[string]string)
		for _, y := range files {
			for _, x := range y["packs"].(bson.A) {
				Names[x.(bson.M)["name"].(string)] = x.(bson.M)["title"].(string)
			}
		}
		return Names
	}
}

func Update_count(user_id int64, name string, _type string) {
	filter := bson.M{"user_id": user_id, "type": _type}
	s := stickers.FindOne(context.TODO(), filter)
	var stick bson.M
	s.Decode(&stick)
	packs := stick["packs"].(bson.A)
	c := packs[len(packs)-1].(bson.M)["count"].(int32)
	c++
	packs[len(packs)-1] = bson.M{"name": packs[len(packs)-1].(bson.M)["name"].(string), "count": c}
	stickers.UpdateOne(context.TODO(), filter, bson.D{{Key: "$set", Value: bson.D{{Key: "packs", Value: packs}}}}, opts)
}
