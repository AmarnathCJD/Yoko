package db

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var approve = database.Collection("approve")

var AllApproved = LoadApproved()

type Approved struct {
	Users []int64 `bson:"users"`
}

func Approve(chatID int64, userID int64) {
	var a Approved
	AP := approve.FindOne(context.TODO(), bson.M{"chat_id": chatID})
	if AP.Err() == nil {
		AP.Decode(&a)
	}
	a.Users = append(a.Users, userID)
	approve.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.M{"$set": bson.M{"users": a.Users}}, opts)
}

func Unapprove(chatID int64, userID int64) {
	var a Approved
	AP := approve.FindOne(context.TODO(), bson.M{"chat_id": chatID})
	if AP.Err() == nil {
		AP.Decode(&a)
	}
	a.Users = Remove(a.Users, userID).([]int64)
	approve.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.M{"$set": bson.M{"users": a.Users}}, opts)
}

func IsApproved(chatID int64, userID int64) bool {
	if Chat, ok := AllApproved[chatID]; ok {
		for _, v := range Chat.Users {
			if v == userID {
				return true
			}
		}
	} else {
		return false
	}
	return false
}

func GetAllApproved(chatID int64) []int64 {
	if Chat, ok := AllApproved[chatID]; ok {
		return Chat.Users
	} else {
		return nil
	}
}

func GetApproved(chatID int64) []int64 {
	var a Approved
	AP := approve.FindOne(context.TODO(), bson.M{"chat_id": chatID})
	if AP.Err() != nil {
		return nil
	}
	AP.Decode(&a)
	return a.Users
}

func LoadApproved() map[int64]Approved {
	var a []bson.M
	var Appr = make(map[int64]Approved)
	AP, _ := approve.Find(context.TODO(), bson.M{})
	if AP.Err() != nil {
		return Appr
	}
	AP.All(context.TODO(), &a)
	for _, v := range a {
		Appr[v["chat_id"].(int64)] = Approved{Users: v["users"].([]int64)}
	}
	return Appr
}

func UnapproveAll(chatID int64) {
	approve.DeleteOne(context.TODO(), bson.M{"chat_id": chatID})
}
