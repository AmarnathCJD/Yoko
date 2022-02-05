package modules


import (
 tb "gopkg.in/tucnak/telebot.v3"
 "fmt"
)

func PARSET(c tb.Context) error {
fmt.Println(ParseMD(c))

}

func ParseMD(c tb.Context) string {
return fmt.Sprint(c.Message().Entites)

}
