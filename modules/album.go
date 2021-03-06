package modules

import (
	"encoding/json"
	"fmt"
	"time"

	tb "gopkg.in/telebot.v3"
)

func OnMediaHandler(c tb.Context) error {
	if afk := AFK(c); afk {
		return nil
	}
	if c.Message().AlbumID != string("") {
		fmt.Println("Album detected")
	}
	return nil
}

func copyAlbum(c tb.Context) {
	if c.Message().Photo != nil {
		AppendAlbum(c.Message().AlbumID, time.Now().Unix(), c.Message().Photo, c)
	}
}

type Album struct {
	AlbumID   string
	TimeStamp int64
	Files     tb.Album
}

var a []Album

func AlbumExist(id string) (bool, int) {
	for i, b := range a {
		if b.AlbumID == id {
			return true, i
		}
	}
	return false, 0
}

func AppendAlbum(data string, _time int64, f *tb.Photo, c tb.Context) {
	e, i := AlbumExist(data)
	if !e {
		a = append(a, Album{data, _time, tb.Album{}})
		a[i].Files = append(a[i].Files, f)
		_, i = AlbumExist(data)
		go func() {
			time.Sleep(time.Second * 2)
			SendAlbum(i, c)
		}()
	} else {
		a[i].Files = append(a[i].Files, f)
	}
}

func SendAlbum(_id int, c tb.Context) {
	msgs := a[_id].Files
	b, _ := json.Marshal(msgs)
	fmt.Println(string(b))
	_, err := c.Bot().SendAlbum(c.Chat(), msgs)
	if err != nil {
		fmt.Println(err)
	}
}
