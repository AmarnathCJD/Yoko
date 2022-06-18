package main

import (
	"fmt"
	"time"

	tb "gopkg.in/telebot.v3"
)

func OnMediaHandler(c tb.Context) error {
	if c.Message().AlbumID != string("") {
		copyAlbum(c)
	}
	return nil
}

func copyAlbum(c tb.Context) {
	AppendAlbum(c.Message().AlbumID, time.Now().Unix(), c.Message().Photo)
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

func AppendAlbum(data string, _time int64, f *tb.Photo) {
	e, i := AlbumExist(data)
	if !e {
		a = append(a, Album{data, _time, tb.Album{}})
		a[i].Files = append(a[i].Files, f)
		_, i = AlbumExist(data)
		go func() {
			time.Sleep(time.Second * 1)
			SendAlbum(i)
		}()
	} else {
		a[i].Files = append(a[i].Files, f)
	}
}

func SendAlbum(_id int) {
	fmt.Println(a[_id])
}
