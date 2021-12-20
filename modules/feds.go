package modules

import (
	"fmt"

	db "github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/tucnak/telebot.v3"
)


func New_fed(c tb.Context) error {
	m := c.Message()
	if !m.Private() {
		c.Reply("Create your federation in my PM - not in a group.")
		return nil
	}
	fed, _, fedname := db.Get_fed_by_owner(m.Sender.ID)
	if fed {
		c.Reply(fmt.Sprintf("You already have a federation called <code>%s</code> ; you can't create another. If you would like to rename it, use <code>/renamefed</code>.", fedname))
		return nil
	}
	if m.Payload == string("") {
		c.Reply("You need to give your federation a name! Federation names can be up to 64 characters long.")
		return nil
	} else if len(m.Payload) > 64 {
		c.Reply("Federation names can only be upto 64 charactors long.")
		return nil
	}
	fed_uid, _ := db.Make_new_fed(m.Sender.ID, m.Payload)
	c.Reply(fmt.Sprintf("Created new federation with FedID: <code>%s</code>.\nUse this ID to join the federation! eg:\n<code>/joinfed %s</code>", fed_uid, fed_uid))
	return nil
}

func Delete_fed(c tb.Context) error {
        m := c.Message()
	if !m.Private() {
		c.Reply("Delete your federation in my PM - not in a group.")
		return nil
	}
        fed, fed_id, fedname := db.Get_fed_by_owner(m.Sender.ID)
        if !fed {
           c.Reply("It doesn't look like you have a federation yet!")
           return nil
        }
        menu := &tb.ReplyMarkup{}
        menu.Inline(menu.Row(menu.Data("Delete Federation", fmt.Sprintf("delfed_%s", fed_id))), menu.Row(menu.Data("Cancel", "cancel_fed_delete")))
        c.Reply(fmt.Sprintf("Are you sure you want to delete your federation? This action cannot be undone - you will lose your entire ban list, and '%s' will be permanently gone.", fedname), menu)
        return nil
}

func Rename_fed(c tb.Context) error {
	m := c.Message()
	if m.Private() {
		c.Reply("You can only rename your fed in PM.")
		return nil
	}
	fed, fed_id, fedname := db.Get_fed_by_owner(m.Sender.ID)
	if !fed {
		c.Reply("It doesn't look like you have a federation yet!")
		return nil
	} else if m.Payload == string("") {
		c.Reply("You need to give your federation a new name! Federation names can be up to 64 characters long.")
		return nil
	}
	db.Rename_fed_by_id(fed_id, m.Payload)
	c.Reply(fmt.Sprintf("Tada! I've renamed your federation from '%s' to '%s'. (FedID: <code>%s</code>).", fedname, m.Payload, fed_id))
        return nil
}
