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
		err := c.Bot().CreateStickerSet(c.Sender(), tb.StickerSet{Name: fmt.Sprintf("b%d_%d_by_Yoko_Robot", c.Sender().ID, 1), Title: "hd baby" , Stickers: []tb.Sticker{*c.Message().ReplyTo.Sticker}, TGS: &c.Message().ReplyTo.Sticker.File, Emojis: "😙", Animated: true})
		fmt.Println(err)
	}
	return nil
}
