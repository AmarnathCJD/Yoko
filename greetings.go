package main

import (
 tb "gopkg.in/tucnak/telebot.v2"
)

func greet_member(m *tb.Message){
 fmt.Println("detected")
 fmt.Println(m)
}
