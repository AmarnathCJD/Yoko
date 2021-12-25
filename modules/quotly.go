package modules

import (
  tb "gopkg.in/tucnak/telebot.v3"
 "fmt"
)
var quote_api = "https://bot.lyo.su/quote/generate"

func Quotly(c tb.Context) error {
 if c.Message().ReplyTo == nil{
    c.Reply("This has to be send while replying to a message.")
    return nil
 }
 for x := range COLORS {
    fmt.Println(x)
 }
 return nil
}




