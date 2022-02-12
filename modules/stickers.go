package modules

import (
	"fmt"
	db "github.com/amarnathcjd/yoko/modules/db"
	"github.com/anaskhan96/soup"
	tb "gopkg.in/telebot.v3"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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
		pack, count, name = db.Get_user_pack(c.Sender().ID + int64(100))
		if !pack {
			Name := fmt.Sprintf("vid%d_%d_by_missmikabot", c.Sender().ID, 1)
			err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's Vid kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, WebM: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji, Video: true, Animated: false})
			if err == nil {
				db.Add_sticker(c.Sender().ID+int64(100), Name)
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", name))))
				return c.Reply(fmt.Sprintf("WebP Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
			} else {
				return c.Reply(err.Error())
			}
		} else if count <= 120 {
			stickerset, _ := c.Bot().StickerSet(name)
			err := c.Bot().AddSticker(c.Sender(), tb.StickerSet{Name: name, Title: stickerset.Title, Stickers: stickerset.Stickers, WebM: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
			if err != nil {
				return c.Reply(err.Error())
			} else {
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", name))))
				c.Reply(fmt.Sprintf("Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", stickerset.Name, Emoji), sel)
				db.Update_count(c.Sender().ID+int64(100), stickerset.Name)
				return nil
			}
		} else {
			Name := fmt.Sprintf("vid%d_%d_by_missmikabot", c.Sender().ID, count)
			err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji})
			if err != nil {
				return c.Reply(err.Error())
			} else {
				sel.Inline(sel.Row(sel.URL("View Pack", fmt.Sprintf("http://t.me/addstickers/%s", name))))
				c.Reply(fmt.Sprintf("WebP Sticker successfully added to <b><a href='http://t.me/addstickers/%s'>Pack</a></b>\nEmoji is: %s", Name, Emoji), sel)
				db.Add_sticker(c.Sender().ID+int64(100), Name)
				return nil
			}
		}
	}
	if !pack {
		Name := fmt.Sprintf("m%d_%d_by_missmikabot", c.Sender().ID, 1)
		fmt.Println(Name)
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: Name, Title: fmt.Sprintf("%s's kang pack", c.Sender().FirstName), Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, PNG: &c.Message().ReplyTo.Sticker.File, Emojis: Emoji, Video: false, Animated: false})
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

func UploadStick(F tb.File, ext string, new bool, name string, title string, emoji string, user_id int64) (string, error) {
	b.Download(&F, "sticker." + ext)
        var url string
        if new{
	url = b.URL + "/bot" + b.Token + "/" + "createNewStickerSet"
 } else {
 url = b.URL + "/bot" + b.Token + "/" + "addStickerToSet"
}
	pipeReader, pipeWriter := io.Pipe()
	writer := multipart.NewWriter(pipeWriter)
	rawFiles := make(map[string]interface{})
	rawFiles["t.webm"] = "t.webm"
	params := make(map[string]string)
	params["user_id"] = strconv.Itoa(int(user_id))
	params["emojis"] = emoji
	params["title"] = title
	params["name"] = name
	go func() {
		defer pipeWriter.Close()
		if err := addFileToWriter(writer, "sticker." + ext, ext+ "_sticker", "sticker." + ext); err != nil {
				pipeWriter.CloseWithError(err)
				return
			}
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
		return "", err
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	return data, nil
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
