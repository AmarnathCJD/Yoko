package modules

import (
	"fmt"
	db "github.com/amarnathcjd/yoko/modules/db"
	"github.com/anaskhan96/soup"
	tb "gopkg.in/tucnak/telebot.v3"
	"io"
	"net/http"
)

func AddSticker(c tb.Context) error {
	pack, count, name := db.Get_user_pack(c.Sender().ID)
	Emoji := "😙"
	if c.Message().Payload != string("") {
		Emoji = c.Message().Payload
	}
	if c.Message().ReplyTo == nil {
		return nil
	}
	if c.Message().ReplyTo.Photo != nil {
		c.Reply("sticker file can only be valid wepb files.")
		return nil
	}
	if c.Message().ReplyTo.Sticker == nil {
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

func CombotSticker(c tb.Context) error {
	query := c.Message().Payload
	req, _ := http.NewRequest("GET", "https://combot.org/telegram/stickers", nil)
	q := req.URL.Query()
	q.Add("q", query)
	req.URL.RawQuery = q.Encode()
	resp, err := myClient.Do(req)
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		c.Reply("Search Service is Down!")
	}
	body, _ := io.ReadAll(resp.Body)
	doc := soup.HTMLParse(string(body))
	x := doc.FindAll("a", "class", "sticker-pack__btn")
	y := doc.FindAll("div", "sticker-pack__title")
	fmt.Println(x, y)
	return nil
}
