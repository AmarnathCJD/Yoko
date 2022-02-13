package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
stickers = database.Collection("stick")
)

type PACK struct {
Name string
Ext string
Title string
}

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

func Get_user_packs(user_id int64) []PACK {
	var s []PACK
        for _, x := range []string{"png", "tgs", "webm"} {
f := bson.M{"user_id": user_id, "type": x}
st := stickers.FindOne(context.TODO(), f)
if st.Err() == nil {
var pk bson.M
st.Decode(&pk)
s = append (s, PACK{pk["name"].(string), x, pk["title"].(string)})
}
}
return s
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
