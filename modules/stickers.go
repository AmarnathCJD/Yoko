package modules

import (
	"fmt"
        db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)

func AddSticker(c tb.Context) error {
	x := db.get_user_pack(c.Sender.ID)
	fmt.Println(x)
        return nil
}
