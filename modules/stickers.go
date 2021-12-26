package modules

import (
	"fmt"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func AddSticker(c tb.Context) error {
	pack, count, name := db.Get_user_pack(c.Sender().ID)
	Emoji := "😙"
	if c.Message().Payload != string("") {
		Emoji = c.Message().Payload
	}
	if c.Message().ReplyTo.Photo != nil {
		c.Reply("sticker file can only be valid wepb files.")
		return nil
	} else if c.Message().ReplyTo.Sticker == nil {
		c.Reply("Yeah, I can't kang that.")
		return nil
	}
	if !pack {
		Name := fmt.Sprintf("d%d_%d_by_Yoko_Robot", c.Sender().ID, 1)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err == nil {
			db.Add_sticker(c.Sender().ID, Name)
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji))
		} else {
			c.Reply(err.Error())
		}
	} else if count <= 120 {
		stickerset, _ := c.Bot().StickerSet(name)
		err := c.Bot().AddSticker(c.Sender(), tb.StickerSet{Name: name, Title: stickerset.Title, Stickers: stickerset.Stickers, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji))
			db.Update_count(c.Sender().ID, stickerset.Name)
		}
	} else {
		Name := fmt.Sprintf("d%d_%d_by_Yoko_Robot", c.Sender().ID, count)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji))
			db.Add_sticker(c.Sender().ID, Name)
		}
	}
	return nil
}

func TGStest(c tb.Context) error {
	x := fmt.Sprintf("cj%d_%d_by_Yoko_Robot", c.Sender().ID, 1)
	file, _ := b.FileByID(c.Message().ReplyTo.Sticker.FileID)
	c.Reply(file)
	err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: x, Title: "smd kang tgs test", Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, TGS: &c.Message().Sticker.File, Emojis: "😭", Animated: true})
	if err != nil {
		c.Reply(err.Error())
	}
	return nil
}
