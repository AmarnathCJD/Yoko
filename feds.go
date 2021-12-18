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
 fed_name := m.Payload
 if fed_name == string("") {
    c.Reply("You need to give your federation a name! Federation names can be up to 64 characters long.")
    return nil
 } else if len(fed_name) > 64 {
    c.Reply("Federation names can only be upto 64 charactors long.")
    return nil
 }
 fed_uid, _ := make_new_fed(m.Sender.ID, fed_name)
 c.Reply(fmt.Sprintf("Created new federation with FedID: <code>%s</code>.\nUse this ID to join the federation! eg:\n<code>/joinfed %s</code>", fed_uid, fed_uid))
 return nil
}
