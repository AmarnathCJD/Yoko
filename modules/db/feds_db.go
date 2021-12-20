package db

import (
	"context"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

var feds = database.Collection("feda")

func Make_new_fed(user_id int64, fedname string) (string, string) {
	uid := uuid.New().String()
	filter := bson.M{"user_id": user_id}
	feds.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"fed_id", uid}, {"fedname", fedname}}}}, opts)
	return uid, fedname
}

func Get_fed_by_owner(user_id int64) (bool, string, string) {
	filter := bson.M{"user_id": user_id}
	fed := feds.FindOne(context.TODO(), filter)
	if fed.Err() != nil {
		return false, "", ""
	}
	var fed_info bson.M
	fed.Decode(&fed_info)
	return true, fed_info["fed_id"].(string), fed_info["fedname"].(string)
}

func Delete_fed(fed_id string) {
	filter := bson.M{"fed_id": fed_id}
	feds.DeleteOne(context.TODO(), filter)
}

func Rename_fed(fed_id string, name string) {
        filter := bson.M{"fed_id": fed_id}
        feds.UpdateOne(context.TODO(), filter, bson.D{{"$set", bson.D{{"fedname", name}}}}, opts)
}
