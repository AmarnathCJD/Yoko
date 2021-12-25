package modules

import (
	"fmt"
        db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func AddSticker(c tb.Context) error {
	pack, _ := db.Get_user_pack(c.Sender().ID)
        fmt.Println(c.Message().Sticker)
	if !pack {
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: "smd", Title: "stfu", Stickers: []*tb.Sticker{c.Message().Sticker}})
		fmt.Println(err)
	}
	return nil
}
