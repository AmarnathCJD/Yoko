package main

import (
 tb "gopkg.in/tucnak/telebot.v3"
 "fmt"
)

func new_fed(c tb.Context) error {
 m := c.Message()
 if !c.Message().Private(){
   b.Reply(m, "Create your federation in my PM - not in a group.")
   return nil
 }
 fed := get_fed_by_owner(m.Sender.ID)
 if fed[0] {
    c.Reply(fmt.Sprintf("You already have a federation called <code>%s</code> ; you can't create another. If you would like to rename it, use <code>/renamefed</code>.", fed[2]))
    return nil
 }
 return nil
}
