package modules

import "fmt"

func AddSticker(c tb.Context) error {
	x := db.get_user_pack(c.Sender.ID)
	fmt.Println(x)
}
