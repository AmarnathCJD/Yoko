package modules

import (
  tb "gopkg.in/tucnak/telebot.v3"
)

var quote_api = "https://bot.lyo.su/quote/generate"

func Quotly(c tb.Context) error {
 if c.Message().ReplyTo == nil{
    c.Reply("This has to be send while replying to a message.")
    return nil
 }
 return nil
}
