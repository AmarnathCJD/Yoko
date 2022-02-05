package modules

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
)

func PARSET(c tb.Context) error {
        fmt.Println(ParseMD(c))
	return nil

}

func ParseMD(c tb.Context) string {
	return fmt.Sprint(c.Message().Entities)

}
