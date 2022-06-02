package modules

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	db "github.com/amarnathcjd/yoko/modules/db"
	"github.com/anaskhan96/soup"
	"github.com/nfnt/resize"
	tb "gopkg.in/telebot.v3"
)

var (
	PACK_TYPES = []string{"png", "webm", "tgs"}
)

func AddSticker(c tb.Context) error {
	Pack, Count := db.GetPack(c.Sender().ID, "png")
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
	if Emoji == string("") {
		Emoji = Sticker.Emoji
		if Emoji == string("") {
			Emoji = "😙"
		}
	}
	if Sticker.Video || Sticker.Animated {
		var Ext = "webm"
		var PrePre = "vidpas"
		var Prefix = "Video"
		if Sticker.Animated {
			Ext = "tgs"
			PrePre = "tgspas"
			Prefix = "Animated"
		}
		Pack, Count = db.GetPack(c.Sender().ID, Ext)
		title := fmt.Sprintf("%s's %s kang pack", c.Sender().FirstName, Prefix)
		if Pack.Name == "" {
			Name := fmt.Sprintf("%s%d_%d_by_missmikabot", PrePre, c.Sender().ID, 1)
			err, er := UploadStick(Sticker.File, Ext, true, Name, title, Emoji, c.Sender().ID)
			if err {
				db.AddSticker(c.Sender().ID, Name, title, Ext)
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
				return c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
			} else {
				return c.Reply("Failed to kang, " + er)
			}
		} else if Pack.Count <= 120 {
			stickerset, _ := c.Bot().StickerSet(Pack.Name)
			err, er := UploadStick(Sticker.File, Ext, false, Pack.Name, stickerset.Title, Emoji, c.Sender().ID)
			if !err {
				return c.Reply("Failed to kang, " + er)
			} else {
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", stickerset.Name))))
				c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji), sel)
				db.UpdateCount(c.Sender().ID, Pack.Type)
				return nil
			}
		} else {
			Name := fmt.Sprintf("%s%d_%d_by_missmikabot", PrePre, c.Sender().ID, Count)
			err, er := UploadStick(Sticker.File, Ext, true, Name, title, Emoji, c.Sender().ID)
			if !err {
				return c.Reply("Failed to kang, " + er)
			} else {
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
				c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
				db.AddSticker(c.Sender().ID, Name, title, Ext)
				return nil
			}
		}

	}
	title := fmt.Sprintf("%s's kang pack", c.Sender().FirstName)
	if Pack.Name == "" {
		Name := fmt.Sprintf("pngdf%d_%d_by_missmikabot", c.Sender().ID, 1)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*Sticker}, PNG: &Sticker.File, Emojis: Emoji, Video: false, Animated: false})
		if err == nil {
			db.AddSticker(c.Sender().ID, Name, title, "png")
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
		} else {
			c.Reply(err.Error())
		}
	} else if Pack.Count <= 120 {
		stickerset, _ := c.Bot().StickerSet(Pack.Name)
		err := c.Bot().AddSticker(c.Sender(), tb.StickerSet{Name: Pack.Name, Title: stickerset.Title, Stickers: stickerset.Stickers, PNG: &Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", stickerset.Name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji), sel)
			db.UpdateCount(c.Sender().ID, "png")
		}
	} else {
		Name := fmt.Sprintf("pngdf%d_%d_by_missmikabot", c.Sender().ID, Count)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*Sticker}, PNG: &Sticker.File, Emojis: Emoji})
		if err != nil {
			c.Reply(err.Error())
		} else {
			sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", Name))))
			c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
			db.AddSticker(c.Sender().ID, Name, title, "png")
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
	resp, err := Client.Do(req)
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
	packs := db.GetUserPacks(c.Sender().ID)
	if len(packs) == 0 {
		c.Reply("You have not created any sticker packs, use <code>/kang</code> to save stickers!")
		return nil
	} else {
		msg := "<b>Here are your kang packs.</b>"
		q := 0
		for _, x := range packs {
			q++
			msg += fmt.Sprintf("\n<b>%d. ~</b> <a href='http://t.me/addstickers/%s'>%s</a>", q, x.Name, x.Title)
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
	resp, err := Client.Post(url, writer.FormDataContentType(), pipeReader)
	if err != nil {
		pipeReader.CloseWithError(err)
		return false, "Error with pipeReader."
	}
	defer resp.Body.Close()
	var Resp mapType
	json.NewDecoder(resp.Body).Decode(&Resp)
	var Error string
	if d, ok := Resp["description"]; ok {

		Error = d.(string)
	}
	return Resp["ok"].(bool), Error
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

func IsSticker(c tb.Context) (bool, *tb.Sticker) {
	r := c.Message().ReplyTo
	if r == nil {
		return false, nil
	}
	if r.Sticker != nil {
		return true, r.Sticker
	}
	if r.Photo != nil {
		return ImageToSticker(r.Photo.File, *c.Bot(), *c.Sender())
	}
	if r.Document != nil {
		if strings.Contains(r.Document.FileName, "webm") {
			return true, &tb.Sticker{File: tb.File{FileID: c.Message().ReplyTo.Document.File.FileID}, Video: true, Animated: false, Emoji: ""}
		}
	}
	return false, nil
}

func ImageToSticker(f tb.File, bot tb.Bot, user tb.User) (bool, *tb.Sticker) {
	err := bot.Download(&f, "sticker.webp")
	if err != nil {
		return false, nil
	}
	file, err := os.Open("sticker.webp")
	if err != nil {
		return false, nil
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return false, nil
	}
	m := resize.Resize(512, 0, img, resize.Lanczos3)
	out, err := os.Create("sticker.png")
	if err != nil {
		return false, nil
	}
	defer out.Close()
	err = jpeg.Encode(out, m, &jpeg.Options{Quality: 100})
	if err != nil {
		return false, nil
	}
	sf, err := bot.UploadSticker(&user, &tb.File{FileLocal: "sticker.png"})
	if err != nil {
		return false, nil
	}
	return true, &tb.Sticker{File: *sf, Video: false, Animated: false, Emoji: ""}
}

// sooon
