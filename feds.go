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
 fed, _, fedname := get_fed_by_owner(m.Sender.ID)
 if fed {
    c.Reply(fmt.Sprintf("You already have a federation called <code>%s</code> ; you can't create another. If you would like to rename it, use <code>/renamefed</code>.", fedname))
    return nil
 }
 return nil
}
