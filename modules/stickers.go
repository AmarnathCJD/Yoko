package modules

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	"github.com/anaskhan96/soup"
	tb "gopkg.in/telebot.v3"
)

var (
	PACK_TYPES = []string{"png", "webm", "tgs"}
)

type KangError struct {
Ok string `json:"result"`
Result bool `json:"result"`
Description string `json:"description"`
Error int `json:"error_code"`
}

func AddSticker(c tb.Context) error {
	pack, count, name := db.Get_user_pack(c.Sender().ID, "png")
	var Emoji string
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
		var IsVid = false
		if c.Message().ReplyTo.Document != nil && strings.HasSuffix(c.Message().ReplyTo.Document.FileName, "webm") {
			IsVid = true
		}
		if !IsVid {
			c.Reply("Yeah, I can't kang that.")
			return nil
		}
	}
	var Sticker *tb.Sticker
	if c.Message().ReplyTo.Sticker == nil {
		Sticker = &tb.Sticker{File: tb.File{FileID: c.Message().ReplyTo.Document.File.FileID}, Video: true, Animated: false, Emoji: ""}
	} else {
		Sticker = c.Message().ReplyTo.Sticker
	}
	fmt.Println(Sticker)
	if Emoji == string("") {
		Emoji = Sticker.Emoji
		if Emoji == string("") {
			Emoji = "😙"
		}
	}
	if Sticker.Video || Sticker.Animated {
		var Ext = "webm"
		var PrePre = "vip"
		var Prefix = "Video"
		if Sticker.Animated {
			Ext = "tgs"
			PrePre = "tgp"
			Prefix = "Animated"
		}
		pack, count, name = db.Get_user_pack(c.Sender().ID, Ext)
		title := fmt.Sprintf("%s's %s kang pack", c.Sender().FirstName, Prefix)
		if !pack {
			Name := fmt.Sprintf("%s%d_%d_by_missmikabot", PrePre, c.Sender().ID, 1)
			err, er := UploadStick(Sticker.File, Ext, true, Name, title, Emoji, c.Sender().ID)
			if err {
				db.Add_sticker(c.Sender().ID, Name, title, Ext)
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
				return c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
			} else {
				return c.Reply("Failed to kang, "+er)
			}
		} else if count <= 120 {
			stickerset, _ := c.Bot().StickerSet(name)
			err, er := UploadStick(Sticker.File, Ext, false, name, stickerset.Title, Emoji, c.Sender().ID)
			if !err {
				return c.Reply("Failed to kang, "+er)
			} else {
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", stickerset.Name))))
				c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji), sel)
				db.Update_count(c.Sender().ID, stickerset.Name, Ext)
				return nil
			}
		} else {
			Name := fmt.Sprintf("%s%d_%d_by_missmikabot", PrePre, c.Sender().ID, count)
			err, er := UploadStick(Sticker.File, Ext, true, Name, title, Emoji, c.Sender().ID)
			if !err {
				return c.Reply("Failed to kang, "+er)
			} else {
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
				c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
				db.Add_sticker(c.Sender().ID, Name, title, Ext)
				return nil
			}
		}

	}
	title := fmt.Sprintf("%s's kang pack", c.Sender().FirstName)
	if !pack {
		Name := fmt.Sprintf("pns%d_%d_by_missmikabot", c.Sender().ID, 1)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*Sticker}, PNG: &Sticker.File, Emojis: Emoji, Video: false, Animated: false})
		if err == nil {
			db.Add_sticker(c.Sender().ID, Name, title, "png")
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
		} else {
			c.Reply(err.Error())
		}
	} else if count <= 120 {
		stickerset, _ := c.Bot().StickerSet(name)
		err := c.Bot().AddSticker(c.Sender(), tb.StickerSet{Name: name, Title: stickerset.Title, Stickers: stickerset.Stickers, PNG: &Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", stickerset.Name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji), sel)
			db.Update_count(c.Sender().ID, stickerset.Name, "png")
		}
	} else {
		Name := fmt.Sprintf("pns%d_%d_by_missmikabot", c.Sender().ID, count)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*Sticker}, PNG: &Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
			db.Add_sticker(c.Sender().ID, Name, title, "png")
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
	packs := db.Get_user_packs(c.Sender().ID)
	if len(packs) == 0 {
		c.Reply("You have not created any sticker packs, use <code>/kang</code> to save stickers!")
		return nil
	} else {
		fmt.Println(packs)
		msg := "<b>Here are your kang packs.</b>"
		q := 0
		for _, x := range packs {
			q++
			Addon := ""
			if x.Ext == "tgs" {
				Addon = "- Animated"
			} else if x.Ext == "webm" {
				Addon = "- Video"
			}
			msg += fmt.Sprintf("\n<b>%d. ~</b> <a href='http://t.me/addstickers/%s'>%s</a> %s", q, x.Name, x.Title, Addon)
		}
		c.Reply(msg)
	}
	return nil
}

func UploadStick(F tb.File, ext string, new bool, name string, title string, emoji string, user_id int64) (bool, string) {
	b.Download(&F, "sticker."+ext)
	var url string
	if new {
		url = b.URL + "/bot" + b.Token + "/" + "createNewStickerSet"
	} else {
		url = b.URL + "/bot" + b.Token + "/" + "addStickerToSet"
	}
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	rawFiles := make(map[string]interface{})
	rawFiles["sticker."+ext] = "sticker." + ext
	params := make(map[string]string)
	params["user_id"] = strconv.Itoa(int(user_id))
	params["emojis"] = emoji
	params["title"] = title
	params["name"] = name
	go func() {
		defer pipeWriter.Close()
		if err := addFileToWriter(writer, "sticker."+ext, ext+"_sticker", "sticker."+ext); err != nil {
			pipeWriter.CloseWithError(err)
			return
		}

		for field, value := range params {
			if err := writer.WriteField(field, value); err != nil {
				pipeWriter.CloseWithError(err)
				return
			}
		}
		if err := writer.Close(); err != nil {
			pipeWriter.CloseWithError(err)
			return
		}
	}()
	resp, err := myClient.Post(url, writer.FormDataContentType(), pipeReader)
	if err != nil {
		pipeReader.CloseWithError(err)
		return false, nil
	}
	defer resp.Body.Close()
	var Resp KangError
	json.NewDecoder(resp.Body).Decode(&d)
	return Resp.Ok, Resp.Error
}

func addFileToWriter(writer *multipart.Writer, filename, field string, file interface{}) error {
	var reader io.Reader
	if r, ok := file.(io.Reader); ok {
		reader = r
	} else if path, ok := file.(string); ok {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		reader = f
	} else {
		return nil
	}

	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, reader)
	return err
}
