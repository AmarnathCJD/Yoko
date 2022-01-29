package modules

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v3"
)

func Test(c tb.Context) error {
	f := c.Message().ReplyTo.Sticker.FileID
	c.Reply(f)
	g := bson.A{f, "sticker"}
	p := GetFile(g, "")

	fmt.Println(p.Send(c.Bot(), c.Chat(), &tb.SendOptions{ReplyTo: c.Message()}))
	return nil
}
