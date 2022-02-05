package modules

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
	"regexp"
	"strings"
)

var HyperLink = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
var Bold = regexp.MustCompile(`\*(.*?)\*`)
var Italic = regexp.MustCompile(`\_(.*?)\_`)
var Strike = regexp.MustCompile(`\~(.*?)\~`)
var Underline = regexp.MustCompile(`\_\_(.*?)\_\_`)
var Spoiler = regexp.MustCompile(`\|\|(.*?)\|\|`)

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
			text = string(text[:offset+cor]) + "<u>" + string(text[offset+cor:offset+cor+length]) + "</u>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == tb.EntityItalic {
			text = string(text[:offset+cor]) + "<i>" + string(text[offset+cor:offset+cor+length]) + "</i>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == tb.EntityStrikethrough {
			text = string(text[:offset+cor]) + "<s>" + string(text[offset+cor:offset+cor+length]) + "</s>" + string(text[offset+cor+length:])
			cor += 7

		} else {
			fmt.Println(x.Type)
			text = string(text[:offset+cor]) + "<tgspoiler>" + string(text[offset+cor:offset+cor+length]) + "</tgspolier>" + string(text[offset+cor+length:])
			cor += 23
		}
	}
	Links := HyperLink.FindAllStringSubmatch(text, -1)
	if Links != nil {
		for _, x := range Links {
			if strings.Contains(x[2], "buttonurl") {
				continue
			}
			text = strings.Replace(text, x[0], fmt.Sprintf("<a href='%s'>%s</a>", x[2], x[1]), -1)
		}
	}
	for _, x := range Bold.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<b>"+x[1]+"</b>", -1)

	}
	for _, x := range Italic.FindAllStringSubmatch(text, -1) {
		if match, _ := regexp.Match(`\_\_(.*?)\_\_`, []byte(x[0])); match {
			continue
		}
		text = strings.Replace(text, x[0], "<i>"+x[1]+"</i>", -1)

	}
	for _, x := range Strike.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<s>"+x[1]+"</s>", -1)

	}
	for _, x := range Underline.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<u>"+x[1]+"</u>", -1)

	}
	for _, x := range Spoiler.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<tgspoiler>"+x[1]+"</tgspoiler>", -1)

	}
	fmt.Println(text)
	return text

}
