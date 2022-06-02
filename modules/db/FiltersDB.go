package db

import (
	"context"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	filters       = database.Collection("filters")
	Filters       = make(map[int64][]MsgDB)
	bracket_regex = regexp.MustCompile(`\(.*\)`)
)

type FilterDocument struct {
	ChatID  int64
	Filters []MsgDB
}

func IsFilterExists(chatID int64, name string) bool {
	if F, ok := Filters[chatID]; ok {
		for _, x := range F {
			if x.Name == name {
				return true
			}
		}
	}
	return false
}

func SaveFilter(chatID int64, msg MsgDB) error {
	if bracket_regex.MatchString(msg.Name) {
		return SaveFilters(chatID, msg)
	}
	_, err := filters.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$push", Value: bson.D{{Key: "filters", Value: msg}}}}, opts)
	if IsFilterExists(chatID, msg.Name) {
		Filters[chatID] = DupFunc(Filters[chatID], msg.Name)
	} else {
		Filters[chatID] = append(Filters[chatID], msg)
	}
	return err
}

func SaveFilters(chatID int64, msg MsgDB) error {
	var names []string
	for _, x := range strings.Split(msg.Name, ",") {
		names = append(names, strings.Trim(x, "()"))
	}
	_, err := filters.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$push", Value: bson.D{{Key: "filters", Value: msg}}}}, opts)
	for _, x := range names {
		if IsFilterExists(chatID, x) {
			Filters[chatID] = DupFunc(Filters[chatID], x)
		} else {
			msg.Name = x
			Filters[chatID] = append(Filters[chatID], msg)
		}
	}
	return err
}

func RemoveFilter(chatID int64, name string) error {
	_, err := filters.UpdateOne(context.TODO(), bson.M{"chat_id": chatID}, bson.D{{Key: "$pull", Value: bson.D{{Key: "filters", Value: bson.M{"name": name}}}}}, opts)
	if IsFilterExists(chatID, name) {
		Filters[chatID] = DupFunc(Filters[chatID], name)
	}
	return err
}

func GetFiltersFromDB(chatID int64) []MsgDB {
	var document FilterDocument
	r := filters.FindOne(context.TODO(), bson.M{"chat_id": chatID})
	r.Decode(&document)
	return document.Filters
}

func GetFilter(chatID int64, name string) *MsgDB {
	if F, ok := Filters[chatID]; ok {
		for _, x := range F {
			if x.Name == name {
				return &x
			}
		}
	}
	return nil
}

func GetFilters(chatID int64) []MsgDB {
	if F, ok := Filters[chatID]; ok {
		return F
	} else {
		return nil
	}
}

func PurgeFilters(chatID int64) error {
	_, err := filters.DeleteOne(context.TODO(), bson.M{"chat_id": chatID})
	if err != nil {
		return err
	}
	delete(Filters, chatID)
	return nil
}

func LoadFilters() map[int64][]MsgDB {
	var documents []FilterDocument
	r, err := filters.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil
	}
	r.All(context.TODO(), &documents)
	for _, x := range documents {
		Filters[x.ChatID] = x.Filters
	}
	return Filters
}

func init() {
	LoadFilters()
}
