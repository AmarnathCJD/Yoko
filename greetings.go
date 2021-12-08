package main

import (
 tb "gopkg.in/tucnak/telebot.v2"
 "fmt"
)

func greet_member(m *tb.Message){
 b.Reply(m, "Hi")
 fmt.Println(m)
}
