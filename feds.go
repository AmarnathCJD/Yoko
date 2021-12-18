package main

import (
 tb "gopkg.in/tucnak/telebot.v3
)

func new_fed(c tb.Context) error {
 if !c.Message().Private(){
   b.Reply(c.Message(), "Create your federation in my PM - not in a group.")
   return nil
 }
 return nil
