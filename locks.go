package main

import (
 tb "gopkg.in/tucnak/telebot.v2"
 "fmt"
)

func lock(m *tb.Message){
 if m.Payload == string(""){
    b.Reply(m, "You haven't specified a type to lock.")
 }
 
}
