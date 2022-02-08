package modules

import (
	"fmt"
	"io"
	"net/http"

	db "github.com/amarnathcjd/yoko/modules/db"
	"github.com/anaskhan96/soup"
	tb "gopkg.in/telebot.v3"
)

func AddSticker(c tb.Context) error {
	pack, count, name := db.Get_user_pack(c.Sender().ID)
	Emoji := "😙"
	if c.Message().Payload != string("") {
		Emoji = c.Message().Payload
	}
	if c.Message().ReplyTo == nil {
		return c.Reply("Reply to a sticker to kang it!")
	}
	if c.Message().ReplyTo.Photo != nil {
		c.Reply("sticker file can only be valid wepb files.")
		return nil
	}
	if c.Message().ReplyTo.Sticker == nil {
		c.Reply("Yeah, I can't kang that.")
		return nil
	}
	if c.Message().ReplyTo.Sticker.Video {
		fmt.Println("Video")
	}
	if !pack {
		Name := fmt.Sprintf("s%d_%d_by_missmikabot", c.Sender().ID, 1)
                fmt.Println(Name)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err == nil {
			db.Add_sticker(c.Sender().ID, Name)
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
		} else {
			c.Reply(err.Error())
		}
	} else if count <= 120 {
		stickerset, _ := c.Bot().StickerSet(name)
		err := c.Bot().AddSticker(c.Sender(), tb.StickerSet{Name: name, Title: stickerset.Title, Stickers: stickerset.Stickers, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji), sel)
			db.Update_count(c.Sender().ID, stickerset.Name)
		}
	} else {
		Name := fmt.Sprintf("d%d_%d_by_missmikabot", c.Sender().ID, count)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
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
	y := doc.FindAll("div", "class", "sticker-pack__title")
	duc, msg := []string{}, ""
	if query == string("") {
		msg = "<b>Trending Packs.</b>"
	} else {
		msg = fmt.Sprintf("Search results for <b>%s</b>", query)
	}
	qt := 0
	for i, b := range x {
		if !stringInSlice(b.Attrs()["href"], duc) {
			qt++
			msg += fmt.Sprintf("\n<b>%d. ~</b> <a href='%s'>%s</a>", qt, b.Attrs()["href"], y[i].Text())
			duc = append(duc, b.Attrs()["href"])
		}
	}
	if qt == 0 {
		c.Reply("No Results found for your query!")
		return nil
	}
	c.Reply(msg)
	return nil
}

func MyPacks(c tb.Context) error {
	pack, _, _ := db.Get_user_pack(c.Sender().ID)
	if !pack {
		c.Reply("You have not created any sticker packs, use <code>/kang</code> to save stickers!")
		return nil
	} else {
		packs := db.Get_user_packs(c.Sender().ID)
		msg := "<b>Here are your kang packs.</b>"
		for i, x := range packs {
			stickerset, _ := c.Bot().StickerSet(x)
			msg += fmt.Sprintf("\n<b>%d. ~</b> <a href='http://t.me/addstickers/%s'>%s</a>", i+1, x, stickerset.Title)
		}
		c.Reply(msg)
	}
	return nil
}
