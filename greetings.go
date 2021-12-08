package main

import (
 tb "gopkg.in/tucnak/telebot.v2"
 "fmt"
)

func greet_member(m *tb.Message){
 fmt.Println("detected")
 fmt.Println(m)
}
