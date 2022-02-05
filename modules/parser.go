package modules

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
)

func PARSET(c tb.Context) error {

	return c.Reply(ParseMD(c))

}

func ParseMD(c tb.Context) string {
	text := c.Message().ReplyTo.Text
        cor := 0
	for _, x := range c.Message().ReplyTo.Entities {
		if x.Type == tb.EntityBold {
			offset, length := x.Offset, x.Length
			text = string(text[:offset+cor]) + "<b>" + string(text[offset+cor:offset+cor+length]) + "</b>" + string(text[offset+cor+length:])
                        cor += 6
		} else if 
	}
	fmt.Println(text)
	return text

}
