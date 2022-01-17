package modules

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v3"
)

func Connect_chat(c tb.Context) error {
	if !c.Message().Private() {
		sel.Inline(sel.Row(sel.URL("Connect to chat", fmt.Sprintf("http://t.me/yoko_robot?start=connect_%d", c.Chat().ID))))
		c.Reply("Tap the following button to connect to this chat in PM", sel)
	} else if c.Message().Payload == string(""){
  c.Reply("I need a chat id to connect to!")
 return nil
} else if !isInt(c.Message().Payload){
c.Reply("I expected a chat id, but this isn't a valid integer")
return nil
}
	return nil
}

func private_connect(c tb.Context) error {
args := strings.SplitN(c.Message().Payload, "_", 2)
	chat_id, _ := strconv.Atoi(args[1])
chat, err := c.Bot().ChatByID(chat_id)
sel.Inline(sel.Row(sel.Data("Admin Commands", "connect_ad_cmd")), sel.Row(sel.Data("User commands", "connect_us_cmd")))
c.Reply(fmt.Sprintf("You have been connected to %s!", chat.Title))
return nil
}
