package modules

import (
"time"
tb "gopkg.in/telebot.v3"
)

func OnMediaHandler(c tb.Context) error {
	if c.Message().AlbumID != string("") {
		copyAlbum(c)
}
	return nil
}


func copyAlbum(c) {

}

type Album struct {
	AlbumID   string
	TimeStamp int64
	Files     []string
}

var a []Album

func main() {
	data := "100"
	_time := time.Now().Unix()
	AppendAlbum(data, _time, "1")
	AppendAlbum(data, _time, "2")
	AppendAlbum(data, _time, "3")
	time.Sleep(time.Second * 1)
	AppendAlbum("200", time.Now().Unix(), "4")
}

func AlbumExist(id string) (bool, int) {
	for i, b := range a {
		if b.AlbumID == id {
			return true, i
		}
	}
	return false, 0
}

func AppendAlbum(data string, _time int64, f string) {
	e, i := AlbumExist(data)
	if !e {
		a = append(a, Album{data, _time, []string{f}})
		_, i = AlbumExist(data)
		go func() {
			time.Sleep(time.Second * 1)
			fmt.Println("Hui")
			SendAlbum(i)
		}()
	} else {
		a[i].Files = append(a[i].Files, f)
	}
}

func SendAlbum(_id int) {
	fmt.Println(a[_id])
}
