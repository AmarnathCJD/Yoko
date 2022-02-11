package modules

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"
)

var (
	HyperLink = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	Bold      = regexp.MustCompile(`\*(.*?)\*`)
	Italic    = regexp.MustCompile(`\_(.*?)\_`)
	Strike    = regexp.MustCompile(`\~(.*?)\~`)
	Underline = regexp.MustCompile(`•(.*?)•`)
	Spoiler   = regexp.MustCompile(`\|\|(.*?)\|\|`)
	Code      = regexp.MustCompile("`(.*?)`")
)

func PARSET(c tb.Context) error {

	return c.Reply(ParseMD(c.Message().ReplyTo))

}

func ParseMD(c *tb.Message) string {
	text := c.Text
	cor := 0
	for _, x := range c.Entities {
		offset, length := x.Offset, x.Length
		fmt.Println(offset, length)
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
		} else if x.Type == "spoiler" {
			text = string(text[:offset+cor]) + "<tg-spoiler>" + string(text[offset+cor:offset+cor+length]) + "</tg-spoiler>" + string(text[offset+cor+length:])
			cor += 25
		}
	}
	for _, x := range HyperLink.FindAllStringSubmatch(text, -1) {
		if strings.Contains(x[2], "buttonurl") {
			continue
		}
		text = strings.Replace(text, x[0], fmt.Sprintf("<a href='%s'>%s</a>", x[2], x[1]), -1)
	}
	for _, x := range Bold.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<b>"+x[1]+"</b>", -1)

	}
	for _, x := range Italic.FindAllStringSubmatch(text, -1) {
		pattern, _ := regexp.Compile(`\_\_(.*?)\_\_`)
		if match := pattern.Match([]byte(x[0])); match {
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
		text = strings.Replace(text, x[0], "<tg-spoiler>"+x[1]+"</tg-spoiler>", -1)

	}
	for _, x := range Code.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<code>"+x[1]+"</code>", -1)
	}
	return text

}

var FILLINGS = []FF{{"{first}", 1}, {"{last}", 2}, {"{username}", 3}, {"{fullname}", 5}, {"{id}", 4}, {"{chatname}", 6}, {"{mention}", 7}}

func ParseString(t string, c tb.Context) (string, bool) {
	q, preview := 0, true
	if strings.Contains(t, "{preview}") {
		preview = false
		t = strings.ReplaceAll(t, "{preview}", "")
	}
	for _, f := range FILLINGS {
		if strings.Contains(t, f.F) {
			q++
			t = strings.ReplaceAll(t, f.F, "%["+strconv.Itoa(f.INDEX)+"]s")

		}

	}
	if strings.Contains(t, "{rules}") {
		t = strings.ReplaceAll(t, "{rules}", "[Rules](buttonurl://rules)")
	}
	first := c.Sender().FirstName
	last := c.Sender().LastName
	fullname := first
	if last != string("") {
		fullname += " " + last
	}
	username := c.Sender().Username
	id := strconv.Itoa(int(c.Sender().ID))
	mention := fmt.Sprintf("<a href='tg://user?id=%s'>%s</a>", id, first)
	if username == string("") {
		username = mention
	}
	chatname := c.Chat().Title
	if q != 0 {
		t = fmt.Sprintf(t, first, last, username, id, fullname, chatname, mention)
	}
	return t, preview

}

func ParseFile() {
}

// Soon
