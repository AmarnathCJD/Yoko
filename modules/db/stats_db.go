package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

var (
	chats       = database.Collection("chats")
	users       = database.Collection("users")
	sudo        = database.Collection("sudo")
	Devs, Sudos = __load_devs()
)

type Chat struct {
	Id    int64  `bson:"_id"`
	Title string `bson:"title"`
}

type User struct {
	Id     int64  `bson:"_id"`
	Name   string `bson:"name"`
	Status string `bson:"status"`
}

func RmUser(s []User, r int64) []User {
	for i, v := range s {
		if v.Id == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func AddChat(chat Chat) {
	chats.InsertOne(context.TODO(), chat)
}

func IsChat(chat_id int64) bool {
	c := chats.FindOne(context.TODO(), Chat{Id: chat_id})
	if c.Err() != nil {
		return false
	} else {
		return true
	}
}

func AddUser(user User) {
	users.InsertOne(context.TODO(), user)
}

func IsUser(user_id int64) bool {
	u := users.FindOne(context.TODO(), User{Id: user_id})
	if u.Err() != nil {
		return false
	} else {
		return true
	}
}

func GetAllChats() []Chat {
	var chat []Chat
	cursor, _ := chats.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &chat)
	return chat
}

func GetAllUsers() []User {
	var user []User
	cursor, _ := users.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &user)
	return user
}

func AddSudo(user_id int64, name string) {
	sudo.UpdateOne(context.TODO(), User{Id: user_id}, bson.M{"$set": User{Id: user_id, Name: name, Status: "sudo"}}, opts)
}

func AddDev(user_id int64, name string) {
	sudo.UpdateOne(context.TODO(), User{Id: user_id}, bson.M{"$set": User{Id: user_id, Name: name, Status: "dev"}}, opts)
}

func RemSudo(user_id int64) bool {
	x, _ := sudo.DeleteOne(context.TODO(), User{Id: user_id})
	if x.DeletedCount == 0 {
		return false
	} else {
		Sudos = RmUser(Sudos, user_id)
		return true
	}
}

func RemDev(user_id int64) bool {
	x, _ := sudo.DeleteOne(context.TODO(), User{Id: user_id})
	if x.DeletedCount == 0 {
		return false
	} else {
		Devs = RmUser(Devs, user_id)
		return true
	}

}

func ListSudo() string {
	var user []User
	var msg = "<b>Bot sudo list.</b>\n"
	cursor, _ := sudo.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &user)
	for _, x := range user {
		if x.Status == "sudo" {
			msg += fmt.Sprintf("<b>~</b> <i>%s</i> (<code>%d</code>)\n", x.Name, x.Id)
		}
	}
	return msg
}

func ListDev() string {
	var user []User
	var msg = "<b>Bot dev list.</b>\n"
	cursor, _ := sudo.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &user)
	for _, x := range user {
		if x.Status == "dev" {
			msg += fmt.Sprintf("<b>~</b><i>%s</i> (<code>%d</code>)\n", x.Name, x.Id)
		}
	}
	return msg
}

func __load_devs() ([]User, []User) {
	var user []User
	var d []User
	var s []User
	cursor, _ := sudo.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &user)
	for _, x := range user {
		if x.Status == "dev" {
			d = append(d, x)
		} else if x.Status == "sudo" {
			s = append(s, x)
		}
	}
	return d, s
}
