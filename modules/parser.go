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
		offset, length := x.Offset, x.Length
		if x.Type == tb.EntityBold {
			text = string(text[:offset+cor]) + "<b>" + string(text[offset+cor:offset+cor+length]) + "</b>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == tb.EntityCode {
			text = string(text[:offset+cor]) + "<code>" + string(text[offset+cor:offset+cor+length]) + "</code>" + string(text[offset+cor+length:])
			cor += 13
		} else if x.Type == tb.EntityUnderline {
text = text = string(text[:offset+cor]) + "<code>" + string(text[offset+cor:offset+cor+length]) + "</code>" + string(text[offset+cor+length:])
			cor += 7
} else if x.Type == tb.EntityItalic {
text = string(text[:offset+cor]) + "<code>" + string(text[offset+cor:offset+cor+length]) + "</code>" + string(text[offset+cor+length:])
			cor += 7
}
	}
	fmt.Println(text)
	return text

}
