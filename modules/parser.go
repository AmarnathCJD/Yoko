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
        for i, x := range c.Message().Entites{
if x.Type == tb.EntityBold {
offset, length := x.Offset, x.Length
 text = string(text[:offset]) + "<b>" + string(text[offset:offset+length]) + "</b>" + string(text[offset+length:])
}
}
	return fmt.Sprint(text)

}
