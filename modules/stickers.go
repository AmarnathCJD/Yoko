package modules

import (
	"fmt"
        db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func AddSticker(c tb.Context) error {
	pack, _ := db.Get_user_pack(c.Sender().ID)
        Emoji := "😙"
        if c.Message().Payload != string(""){
           Emoji = c.Message().Payload
        }
	if !pack {
                Name := fmt.Sprintf("c%d_%d_by_Yoko_Robot", c.Sender().ID, 1)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName) , Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err == nil {
                   db.Add_sticker(c.Sender().ID, Name)
                   c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji))
                } else {
                   c.Reply(err.Error())
	}
	return nil
}
