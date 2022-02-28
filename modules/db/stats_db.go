package db

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Gban struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	Reason    string `json:"reason"`
	Admin     int64  `json:"admin"`
	Date      int64  `json:"date"`
}

var (
	chats       = database.Collection("chats")
	users       = database.Collection("users")
	sudo        = database.Collection("sudo")
	gbans       = database.Collection("gbans")
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
	Sudos = append(Sudos, User{Id: user_id, Name: name, Status: "sudo"})
}

func AddDev(user_id int64, name string) {
	sudo.UpdateOne(context.TODO(), User{Id: user_id}, bson.M{"$set": User{Id: user_id, Name: name, Status: "dev"}}, opts)
	Devs = append(Devs, User{Id: user_id, Name: name, Status: "dev"})
}

func RemSudo(user_id int64) bool {
	x, _ := sudo.DeleteOne(context.TODO(), bson.M{"_id": user_id})
	if x.DeletedCount == 0 {
		return false
	} else {
		Sudos = RmUser(Sudos, user_id)
		return true
	}
}

func RemDev(user_id int64) bool {
	x, _ := sudo.DeleteOne(context.TODO(), bson.M{"_id": user_id})
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

func GbanUser(user_id int64, name string, reason string, banner int64) bool {
	if gbans.FindOne(context.TODO(), Gban{Id: user_id}).Err() != nil {
		gbans.InsertOne(context.TODO(), Gban{Id: user_id, FirstName: name, Reason: reason, Admin: banner, Date: time.Now().Unix()})
		return true
	} else {
		gbans.UpdateOne(context.TODO(), Gban{Id: user_id}, bson.M{"$set": Gban{Id: user_id, FirstName: name, Reason: reason, Admin: banner, Date: time.Now().Unix()}}, opts)
		return false
	}
}

func UngbanUser(user_id int64) bool {
	if gbans.FindOne(context.TODO(), Gban{Id: user_id}).Err() != nil {
		return false
	} else {
		gbans.DeleteOne(context.TODO(), Gban{Id: user_id})
		return true
	}
}

func IsGbanned(user_id int64) bool {
	if gbans.FindOne(context.TODO(), Gban{Id: user_id}).Err() != nil {
		return false
	} else {
		return true
	}
}

func GetAllGbans() []Gban {
	var gban []Gban
	cursor, _ := gbans.Find(context.TODO(), bson.M{})
	cursor.All(context.TODO(), &gban)
	return gban
}

func GatherStats() string {
	Stats := "Mika (V1.3.2) Stats\n"
	Stats += fmt.Sprintf("<b>•</b> Database structure version <code>%s</code>", "15")
	Stats += fmt.Sprintf("<b>•</b> <code>%d</code> total users, in <code>%d</code> chats\n", len(GetAllUsers()), len(GetAllChats()))
	return Stats
}
