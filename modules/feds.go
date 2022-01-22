package modules

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	db "github.com/amarnathcjd/yoko/modules/db"
	"go.mongodb.org/mongo-driver/bson"
	tb "gopkg.in/tucnak/telebot.v3"
)

var (
	sel               = &tb.ReplyMarkup{}
	accept_fpromote   = sel.Data("Accept", "accept_fpromote")
	deny_fpromote     = sel.Data("Decline", "decline_fpromote")
	accept_ftransfer  = sel.Data("Accept", "accept_ftransfer")
	deny_ftransfer    = sel.Data("Decline", "deny_ftransfer")
	confirm_ftransfer = sel.Data("Confirm", "confirm_ftransfer")
	reject_ftransfer  = sel.Data("Cancel", "reject_ftransfer")
	check_fed_admins  = sel.Data("Check Fed Admins", "check_fed_admins")
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
	if !m.Private() {
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

func Join_fed(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("Only supergroups can join feds.")
		return nil
	}
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == "administrator" {
		c.Reply("You need to be the chat creator to do this!")
		return nil
	} else if p.Role == tb.Member {
		c.Reply("You need to be an admin to do this!")
		return nil
	} else if p.Role != tb.Creator {
		c.Reply("You need to be chat creator to do this!")
		return nil
	}
	args := c.Message().Payload
	if args == string("") {
		c.Reply("You need to specify which federation you're asking about by giving me a FedID!")
		return nil
	} else if len(args) < 10 {
		c.Reply("This isn't a valid FedID format!")
		return nil
	}
	getfed := db.Search_fed_by_id(args)
	if getfed == nil {
		c.Reply("This FedID does not refer to an existing federation.")
		return nil
	} else {
		c.Reply(fmt.Sprintf("Successfully joined the '%s' federation! All new federation bans will now also remove the members from this chat.", getfed["fedname"]))
		db.Chat_join_fed(args, c.Chat().ID)
	}
	return nil
}

func Leave_fed(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("Only supergroups can join feds.")
		return nil
	}
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role == "administrator" {
		c.Reply("You need to be the chat creator to do this!")
		return nil
	} else if p.Role == tb.Member {
		c.Reply("You need to be an admin to do this!")
		return nil
	} else if p.Role != tb.Creator {
		c.Reply("You need to be chat creator to do this!")
		return nil
	}
	chat_fed := db.Get_chat_fed(c.Chat().ID)
	if chat_fed == string("") {
		c.Reply("This chat isn't currently in any federations!")
	} else {
		fed := db.Search_fed_by_id(chat_fed)
		c.Reply(fmt.Sprintf("Chat %s has left the '%s' federation.", c.Chat().Title, fed["fedname"].(string)))
		db.Chat_leave_fed(c.Chat().ID)
	}
	return nil
}

func Chat_fed(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is for supergroups only!")
		return nil
	}
	p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
	if p.Role != tb.Creator && p.Role != tb.Administrator {
		c.Reply("You need to be an admin to do this.")
		return nil
	}
	fed_id := db.Get_chat_fed(c.Chat().ID)
	if fed_id == string("") {
		c.Reply("This chat isn't part of any feds yet!")
	} else {
		fd := db.Search_fed_by_id(fed_id)
		c.Reply(fmt.Sprintf("Chat %s is part of the following federation: %s (ID: <code>%s</code>)", c.Chat().Title, fd["fedname"].(string), fed_id))
	}
	return nil
}

func Fpromote(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats, not in pm!")
		return nil
	}
	user, _ := get_user(c.Message())
	if user == nil {
		return nil
	}
	fed, fed_id, fedname := db.Get_fed_by_owner(c.Sender().ID)
	if !fed {
		c.Reply("Only federation creators can promote people, and you don't seem to have a federation to promote to!")
		return nil
	} else if user.ID == c.Sender().ID {
		c.Reply("Yeah well you are the fed owner!")
		return nil
	}
	fban, reason := db.Is_Fbanned(user.ID, fed_id)
	if fban {
		if reason != string("") {
			reason = fmt.Sprintf("\nReason: %s", reason)
		}
		c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> is fbanned in %s. You should unfban them before promoting.%s", user.ID, user.FirstName, fedname, reason))
		return nil
	}
	if db.Is_user_fed_admin(user.ID, fed_id) {
		c.Reply(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> is already an admin in %s!", user.ID, user.FirstName, fedname))
		return nil
	}
	accept_fpromote.Data = strconv.Itoa(int(c.Sender().ID)) + "|" + strconv.Itoa(int(user.ID))
	deny_fpromote.Data = strconv.Itoa(int(c.Sender().ID)) + "|" + strconv.Itoa(int(user.ID))
	sel.Inline(sel.Row(accept_fpromote, deny_fpromote))
	c.Reply(fmt.Sprintf("Please get <a href='tg://user?id=%d'>%s</a> to confirm that they would like to be fed admin for %s", user.ID, user.FirstName, fedname), sel)
	return nil
}

func Fpromote_cb(c tb.Context) error {
	data := strings.SplitN(c.Callback().Data, "|", 2)
	owner_id_int, _ := strconv.Atoi(data[0])
	user_id_int, _ := strconv.Atoi(data[1])
	owner_id := int64(owner_id_int)
	user_id := int64(user_id_int)
	if c.Sender().ID != user_id {
		c.Respond(&tb.CallbackResponse{Text: "You are not the user being fpromoted", ShowAlert: true})
		return nil
	}
	_, fed_id, fedname := db.Get_fed_by_owner(owner_id)
	db.User_join_fed(fed_id, user_id)
	c.Edit(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> is now an admin of %s (<code>%s</code>)", user_id, c.Sender().FirstName, fedname, fed_id))
	return nil
}

func Fpromote_deny_cb(c tb.Context) error {
	data := strings.SplitN(c.Callback().Data, "|", 2)
	owner_id_int, _ := strconv.Atoi(data[0])
	user_id_int, _ := strconv.Atoi(data[1])
	owner_id := int64(owner_id_int)
	user_id := int64(user_id_int)
	if c.Sender().ID == owner_id {
		c.Edit(fmt.Sprintf("Fedadmin promotion cancelled by <a href='tg://user?id=%d'>%s</a>", c.Sender().ID, c.Sender().FirstName))
		return nil
	} else if c.Sender().ID == user_id {
		c.Edit(fmt.Sprintf("Fedadmin promotion has been refused by <a href='tg://user?id=%d'>%s</a>", c.Sender().ID, c.Sender().FirstName))
		return nil
	} else {
		c.Respond(&tb.CallbackResponse{Text: "You are not the user being fpromoted", ShowAlert: true})
		return nil
	}
}

func Fdemote(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats, not in pm!")
		return nil
	}
	user, _ := get_user(c.Message())
	if user == nil {
		return nil
	}
	fed, fed_id, fedname := db.Get_fed_by_owner(c.Sender().ID)
	if !fed {
		c.Reply("Only federation creators can demote people, and you don't seem to have a federation to promote to!")
		return nil
	} else if user.ID == c.Sender().ID {
		c.Reply("Yeah well you are the fed owner!")
		return nil
	}
	if !db.Is_user_fed_admin(user.ID, fed_id) {
		c.Reply(fmt.Sprintf("This person isn't a federation admin for '%s', how could I demote them?", fedname))
		return nil
	}
	db.User_leave_fed(fed_id, user.ID)
	c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> is no longer an admin of %s (<code>%s</code>)", user.ID, user.FirstName, fedname, fed_id))
	return nil
}

func Transfer_fed_user(c tb.Context) error {
	if c.Message().Private() {
		c.Reply("This command is made to be used in group chats, not in pm!")
		return nil
	}
	user, _ := get_user(c.Message())
	if user == nil {
		return nil
	}
	fed, fed_id, fedname := db.Get_fed_by_owner(c.Sender().ID)
	if !fed {
		c.Reply("You don't have a fed to transfer!")
		return nil

	} else if user.ID == c.Sender().ID {
		c.Reply("You can only transfer your fed to others!")
		return nil
	}
	fed_2, _, _ := db.Get_fed_by_owner(user.ID)
	if fed_2 {
		c.Reply(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> already owns a federation - they can't own another.", user.ID, user.FirstName))
		return nil
	}
	if !db.Is_user_fed_admin(user.ID, fed_id) {
		c.Reply(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> isn't an admin in %s - you can only give your fed to other admins.", user.ID, user.FirstName, fedname))
		return nil
	}
	accept_ftransfer.Data = strconv.Itoa(int(c.Sender().ID)) + "|" + strconv.Itoa(int(user.ID))
	deny_ftransfer.Data = strconv.Itoa(int(c.Sender().ID)) + "|" + strconv.Itoa(int(user.ID))
	sel.Inline(sel.Row(accept_ftransfer, deny_ftransfer))
	c.Reply(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>, please confirm you would like to receive fed %s (<code>%s</code>) from <a href='tg://user?id=%d'>%s</a>", user.ID, user.FirstName, fedname, fed_id, c.Sender().ID, c.Sender().FirstName), sel)
	return nil
}

func Accept_Transfer_fed_cb(c tb.Context) error {
	data := strings.SplitN(c.Callback().Data, "|", 2)
	owner_id_int, _ := strconv.Atoi(data[0])
	user_id_int, _ := strconv.Atoi(data[1])
	owner_id := int64(owner_id_int)
	user_id := int64(user_id_int)
	if c.Sender().ID != user_id {
		c.Respond(&tb.CallbackResponse{Text: "This action is not intended for you.", ShowAlert: true})
		return nil
	}
	_, fed_id, fedname := db.Get_fed_by_owner(owner_id)
	owner, _ := c.Bot().ChatByID(owner_id)
	confirm_ftransfer.Data = strconv.Itoa(int(owner_id)) + "|" + strconv.Itoa(int(c.Sender().ID))
	reject_ftransfer.Data = strconv.Itoa(int(owner_id)) + "|" + strconv.Itoa(int(c.Sender().ID))
	sel.Inline(sel.Row(confirm_ftransfer, reject_ftransfer))
	c.Edit(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>, please confirm that you wish to send fed %s (<code>%s</code>) to <a href='tg://user?id=%d'>%s</a> this cannot be undone.", owner_id, owner.FirstName, fedname, fed_id, c.Sender().ID, c.Sender().FirstName), sel)
	return nil
}

func Decline_Transfer_fed_cb(c tb.Context) error {
	data := strings.SplitN(c.Callback().Data, "|", 2)
	owner_id_int, _ := strconv.Atoi(data[0])
	user_id_int, _ := strconv.Atoi(data[1])
	owner_id := int64(owner_id_int)
	user_id := int64(user_id_int)
	if c.Sender().ID != user_id && c.Sender().ID != owner_id {
		c.Respond(&tb.CallbackResponse{Text: "This action is not intended for you.", ShowAlert: true})
		return nil
	}
	if c.Sender().ID == user_id {
		c.Edit(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> has declined the fed transfer.", c.Sender().ID, c.Sender().FirstName))
	} else if c.Sender().ID == owner_id {
		c.Edit(fmt.Sprintf("<a href='tg://user?id=%d'>%s</a> has cancelled the fed transfer.", c.Sender().ID, c.Sender().FirstName))
	}
	return nil
}

func Confirm_Transfer_Fed_cb(c tb.Context) error {
	data := strings.SplitN(c.Callback().Data, "|", 2)
	owner_id_int, _ := strconv.Atoi(data[0])
	user_id_int, _ := strconv.Atoi(data[1])
	owner_id := int64(owner_id_int)
	user_id := int64(user_id_int)
	if c.Sender().ID != owner_id {
		c.Respond(&tb.CallbackResponse{Text: "This action is not intended for you.", ShowAlert: true})
		return nil
	}
	user, _ := c.Bot().ChatByID(user_id)
	_, fed_id, fedname := db.Get_fed_by_owner(owner_id)
	c.Edit(fmt.Sprintf("Congratulations! Federation %s (<code>%s</code>) has successfully been transferred from <a href='tg://user?id=%d'>%s</a> to <a href='tg://user?id=%d'>%s</a>.", fedname, fed_id, c.Sender().ID, c.Sender().FirstName, user.ID, user.FirstName))
	db.Transfer_fed(fed_id, user_id)
	db.User_leave_fed(fed_id, user_id)
	c.Send(fmt.Sprintf("<b>Fed Transfer</b>\n<b>Fed:</b> %s\n<b>New Fed Owner:</b> <a href='tg://user?id=%d'>%s</a> - <code>%d</code>\n<b>Old Fed Owner:</b> <a href='tg://user?id=%d'>%s</a> - <code>%d</code>\n<a href='tg://user?id=%d'>%s</a> is now the fed owner. They can promote/demote admins as they like.", fedname, user.ID, user.FirstName, user.ID, c.Sender().ID, c.Sender().FirstName, c.Sender().ID, user.ID, user.FirstName))
	return nil
}

func Deny_Transfer_Fed_cb(c tb.Context) error {
	data := strings.SplitN(c.Callback().Data, "|", 2)
	owner_id_int, _ := strconv.Atoi(data[0])
	owner_id := int64(owner_id_int)
	if c.Sender().ID != owner_id {
		c.Respond(&tb.CallbackResponse{Text: "This action is not intended for you.", ShowAlert: true})
		return nil
	}
	c.Edit(fmt.Sprintf("Fed transfer has been cancelled by <a href='tg://user?id=%d'>%s</a>.", c.Sender().ID, c.Sender().FirstName))
	return nil
}

func Fban(c tb.Context) error {
	fed_id, have_fed, fban_msg := "", false, ""
	if !c.Message().Private() {
		fed_id = db.Get_chat_fed(c.Chat().ID)
		if fed_id == string("") {
			c.Reply("This chat isn't in any federations.")
			return nil
		}
		if !db.Is_user_fed_admin(c.Sender().ID, fed_id) {
			fed := db.Search_fed_by_id(fed_id)
			c.Reply(fmt.Sprintf("You aren't a federation admin for %s!", fed["fedname"].(string)))
			return nil
		}
	} else {
		have_fed, fed_id, _ = db.Get_fed_by_owner(c.Sender().ID)
		if !have_fed {
			c.Reply("You aren't the creator of any feds to act in.")
			return nil
		}
	}
	u, x := get_user(c.Message())
	if u == nil {
		return nil
	}
	if len(x) > 1024 {
		x = x[:1024] + "\n\nNote: The fban reason was over 1024 characters, so has been truncated."
	}
	fed := db.Search_fed_by_id(fed_id)
	fedname := fed["fedname"].(string)
	if u.ID == BOT_ID {
		c.Reply("Oh you're a funny one aren't you! I am not going to fedban myself.")
		return nil
	} else if IS_SUDO(u.ID) {
		c.Reply("I'm not banning one of my sudo users.")
		return nil
	} else if db.Is_user_fed_admin(u.ID, fed_id) {
		c.Reply(fmt.Sprintf("I'm not banning a fed admin/owner from their own fed! (%s)", fedname))
		return nil
	}
	is_banned, r := db.Is_Fbanned(u.ID, fed_id)
	if is_banned && x == string("") && r == string("") {
		c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> is already banned in %s. There is no reason set for their fedban yet, so feel free to set one.", u.ID, u.FirstName, fedname))
		return nil
	} else if is_banned && x == r {
		c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> has already been fbanned, with the exact same reason.", u.ID, u.FirstName))
		return nil
	} else if is_banned && x == string("") {
		if r == string("") {
			c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> is already banned in %s.", u.ID, u.FirstName, fedname))
		} else {
			c.Reply(fmt.Sprintf("User <a href='tg://user?id=%d'>%s</a> is already banned in %s, with reason:\n<code>%s</code>.", u.ID, u.FirstName, fedname, r))
		}
		return nil
	}
	db.Fban_user(u.ID, fed_id, x, u.FirstName, time.Now().Unix(), c.Sender().ID)
	if !is_banned {
		fban_msg = fmt.Sprintf("<b>New FedBan</b>\n<b>Fed:</b> %s\n<b>FedAdmin:</b> <a href='tg://user?id=%d'>%s</a>\n<b>User:</b> <a href='tg://user?id=%d'>%s</a>\n<b>User ID:</b> <code>%d</code>", fedname, c.Sender().ID, c.Sender().FirstName, u.ID, u.FirstName, u.ID)
		if x != string("") {
			fban_msg += fmt.Sprintf("\n<b>Reason:</b> %s", x)
		}
	} else {
		fban_msg = fmt.Sprintf("<b>FedBan Reason Update</b>\n<b>Fed:</b> %s\n<b>FedAdmin:</b> <a href='tg://user?id=%d'>%s</a>\n<b>User:</b> <a href='tg://user?id=%d'>%s</a>\n<b>User ID:</b> <code>%d</code>", fedname, c.Sender().ID, c.Sender().FirstName, u.ID, u.FirstName, u.ID)
		if r != string("") {
			fban_msg += fmt.Sprintf("\n<b>Previous Reason:</b> %s", r)
		}
		if x != string("") {
			fban_msg += fmt.Sprintf("\n<b>Reason:</b> %s", x)
		}
	}
	c.Send(fban_msg)
	getfednotif := db.Get_FEdnotif(fed_id)
	if getfednotif {
		if c.Chat().ID != fed["user_id"].(int64) {
			c.Bot().Send(&tb.User{ID: fed["user_id"].(int64)}, fban_msg)
		}
	}
	fedchats := fed["chats"].(bson.A)
	if len(fedchats) != 0 {
		for _, x := range fedchats {
			c.Bot().Ban(&tb.Chat{ID: x.(int64)}, &tb.ChatMember{User: u})
		}
	}
	subs := db.Get_fed_subs(fed_id)
	if len(subs) != 0 {
		for _, xd := range subs {
			db.Fban_user(u.ID, xd.(string), x, u.FirstName, time.Now().Unix(), c.Sender().ID)
		}
	}
	return nil
}

func Unfban(c tb.Context) error {
	fed_id, have_fed := "", false
	if !c.Message().Private() {
		fed_id = db.Get_chat_fed(c.Chat().ID)
		if fed_id == string("") {
			c.Reply("This chat isn't in any federations.")
			return nil
		}
		if !db.Is_user_fed_admin(c.Sender().ID, fed_id) {
			fed := db.Search_fed_by_id(fed_id)
			c.Reply(fmt.Sprintf("You aren't a federation admin for %s!", fed["fedname"].(string)))
			return nil
		}
	} else {
		have_fed, fed_id, _ = db.Get_fed_by_owner(c.Sender().ID)
		if !have_fed {
			c.Reply("You aren't the creator of any feds to act in.")
			return nil
		}
	}
	u, x := get_user(c.Message())
	if u == nil {
		return nil
	}
	if len(x) > 1024 {
		x = x[:1024] + "\n\nNote: The unfban reason was over 1024 characters, so has been truncated."
	}
	fed := db.Search_fed_by_id(fed_id)
	fedname := fed["fedname"].(string)
	if u.ID == BOT_ID {
		c.Reply("Oh you're a funny one aren't you! How do you think I would have fbanned myself hm?.")
		return nil
	} else if IS_SUDO(u.ID) {
		c.Reply("I'm not banning one of my sudo users.")
		return nil
	} else if db.Is_user_fed_admin(u.ID, fed_id) {
		c.Reply("fed admin/owner cant be banned!")
		return nil
	}
	is_banned, _ := db.Is_Fbanned(u.ID, fed_id)
	if !is_banned {
		c.Reply(fmt.Sprintf("This user isn't banned in the current federation, %s. (<code>%s</code>)", fed_id, fedname))
		return nil
	}
	unfban_msg := fmt.Sprintf("<b>New un-FedBan</b>\n<b>Fed:</b> %s\n<b>FedAdmin:</b> <a href='tg://user?id=%d'>%s</a>\n<b>User:</b> <a href='tg://user?id=%d'>%s</a>\n<b>User ID:</b> <code>%d</code>", fedname, c.Sender().ID, c.Sender().FirstName, u.ID, u.FirstName, u.ID)
	if x != string("") {
		unfban_msg += fmt.Sprintf("\n<b>Reason:</b> %s", x)
	}
	c.Send(unfban_msg)
	getfednotif := db.Get_FEdnotif(fed_id)
	if getfednotif {
		if c.Chat().ID != fed["user_id"].(int64) {
			c.Bot().Send(&tb.User{ID: fed["user_id"].(int64)}, unfban_msg)
		}
	}
	return nil
}

func sub_fed(c tb.Context) error {
	f, fed_id, fedname := db.Get_fed_by_owner(c.Sender().ID)
	if !f {
		c.Reply("Only federation creators can subscribe to a fed. But you don't have a federation!")
		return nil
	}
	if c.Message().Payload == string("") {
		c.Reply("You need to specify which federation you're asking about by giving me a FedID!")
		return nil
	} else if len(c.Message().Payload) < 10 {
		c.Reply("This isn't a valid FedID format!")
		return nil
	}
	fed := db.Search_fed_by_id(c.Message().Payload)
	if fed == nil {
		c.Reply("This FedID does not refer to an existing federation.")
		return nil
	}
	if c.Message().Payload == fed_id {
		c.Reply("... What's the point in subscribing a fed to itself?")
		return nil
	}
	if len(db.Get_my_subs(fed_id)) >= 5 {
		c.Reply("You can subscribe to at most 5 federations. Please unsubscribe from other federations before adding more.")
		return nil
	}
	c.Reply(fmt.Sprintf("Federation <code>%s</code> has now subscribed to <code>%s</code>. All fedbans in <code>%s</code> will now take effect in both feds.", fedname, fed["fedname"].(string), fed["fedname"].(string)))
	db.SUB_fed(fed["fed_id"].(string), fed_id)
	return nil
}

func unsub_fed(c tb.Context) error {
	f, fed_id, fedname := db.Get_fed_by_owner(c.Sender().ID)
	if !f {
		c.Reply("Only federation creators can unsubscribe to a fed. But you don't have a federation!")
		return nil
	}
	if c.Message().Payload == string("") {
		c.Reply("You need to specify which federation you're asking about by giving me a FedID!")
		return nil
	} else if len(c.Message().Payload) < 10 {
		c.Reply("This isn't a valid FedID format!")
		return nil
	}
	fed := db.Search_fed_by_id(c.Message().Payload)
	if fed == nil {
		c.Reply("This FedID does not refer to an existing federation.")
		return nil
	}
	c.Reply(fmt.Sprintf("Federation <code>%s</code> is no longer subscribed to <code>%s</code>. Bans in <code>%s</code> will no longer be applied. Please note that any bans that happened because the user was banned from the subfed will need to be removed manually.", fedname, fed["fedname"].(string), fed["fedname"].(string)))
	db.UNSUB_fed(fed["fed_id"].(string), fed_id)
	return nil
}

func fed_info(c tb.Context) error {
	if !c.Message().Private() {
		b, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
		check(err)
		if b.Role == tb.Member {
			c.Reply("This command can only be used in private.")
			return nil
		}
	}
	input := c.Message().Payload
	var fed_id string
	if input != string("") {
		if len(input) < 10 {
			c.Reply("This isn't a valid FedID format!")
			return nil
		}
		getfed := db.Search_fed_by_id(input)
		if getfed == nil {
			c.Reply("This FedID does not refer to an existing federation.")
			return nil
		}
		fed_id = getfed["fed_id"].(string)
	} else {
		f, fedid, _ := db.Get_fed_by_owner(c.Sender().ID)
		if !f {
			c.Reply("You need to give me a FedID to check, or be a federation creator to use this command!")
			return nil
		}
		fed_id = fedid
	}
	info := db.Search_fed_by_id(fed_id)
	fadmins, fbans, fchats, subbed, mysubs := len(info["fadmins"].(bson.A)), db.Get_len_fbans(fed_id), len(info["chats"].(bson.A)), db.Get_fed_subs(fed_id), db.Get_my_subs(fed_id)
	F_MSG := fmt.Sprintf("Fed info:\nFedID: <code>%s</code>\nName: %s\nCreator: <a href='tg://user?id=%d'>this person</a> (<code>%d</code>)\nNumber of admins: <code>%d</code>\nNumber of bans: <code>%d</code>\nNumber of connected chats: <code>%d</code>\nNumber of subscribed feds: <code>%d</code>", fed_id, info["fedname"].(string), info["user_id"].(int64), info["user_id"].(int64), fadmins, fbans, fchats, len(subbed))
	if len(mysubs) == 0 {
		F_MSG += "\n\nThis federation is not subscribed to any other feds."
	} else {
		F_MSG += "\n\nSubscribed to the following feds:"
		for _, f := range mysubs {
			fd := db.Search_fed_by_id(f.(string))
			if fd != nil {
				F_MSG += fmt.Sprintf("\n- %s (<code>%s</code>)", fd["fedname"].(string), f)
			}

		}
	}
	check_fed_admins.Data = fed_id
	sel.Inline(sel.Row(check_fed_admins))
	c.Reply(F_MSG, sel)
	return nil
}

func Check_f_admins_cb(c tb.Context) error {
	if !c.Message().Private() {
		b, err := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
		check(err)
		if b.Role == tb.Member {
			c.Respond(&tb.CallbackResponse{Text: "This command can only be used in private.", ShowAlert: true})
			return nil
		}
	}
	fed_id := c.Callback().Data
	fed := db.Search_fed_by_id(fed_id)
	fadmins := db.Get_all_fed_admins(fed_id)
	out_str := fmt.Sprintf("Admins in federation '%s':", fed["fedname"].(string))
	fadmins = append(fadmins, fed["user_id"].(int64))
	for _, x := range fadmins {
		u, err := c.Bot().ChatByID(x.(int64))
		name := "User"
		if err == nil {
			name = u.FirstName
		}
		out_str += fmt.Sprintf("\n- <a href='tg://user?id=%d'>%s</a> (<code>%d</code>)", x.(int64), name, x.(int64))
	}
	c.Bot().EditReplyMarkup(c.Message(), nil)
	c.Reply(out_str)
	return nil
}

func Fednotif(c tb.Context) error {
	if !c.Message().Private() {
		c.Reply("This command is made to be used in PM.")
		return nil
	}
	f, fed_id, fname := db.Get_fed_by_owner(c.Sender().ID)
	if !f {
		c.Reply("You aren't the creator of any feds to act in.")
		return nil
	}
	args := c.Message().Payload
	if args == string("") {
		mode := db.Get_FEdnotif(fed_id)
		if mode {
			c.Reply(fmt.Sprintf("The <code>%s</code> fed is currently sending notifications to it's creator when a fed action is performed.", fname))
			return nil
		} else {
			c.Reply(fmt.Sprintf("The <code>%s</code> fed is currently <b>NOT</b> sending notifications to it's creator when a fed action is performed.", fname))
			return nil
		}
	} else if stringInSlice(args, []string{"on", "yes", "enable"}) {
		c.Reply(fmt.Sprintf("The fed silence setting for <code>%s</code> has been updated to: <code>true</code>", fname))
		db.FEdnotif(fed_id, true)
		return nil
	} else if stringInSlice(args, []string{"off", "no", "disable"}) {
		c.Reply(fmt.Sprintf("The fed silence setting for <code>%s</code> has been updated to: <code>false</code>", fname))
		db.FEdnotif(fed_id, false)
		return nil
	} else {
		c.Reply("Your input was not recognised as one of: yes/no/on/off")
		return nil
	}
}

func FedExport(c tb.Context) error {
	fed_id := ""
	if !c.Message().Private() {
		p, _ := c.Bot().ChatMemberOf(c.Chat(), c.Sender())
		if p.Role != tb.Creator && p.Role != tb.Administrator {
			c.Reply("You need to be an admin to do this.")
			return nil
		}
		fed_id = db.Get_chat_fed(c.Chat().ID)
		if fed_id == string("") {
			c.Reply("This chat isn't in any federations.")
			return nil
		}

	} else {
		fed, fedid, _ := db.Get_fed_by_owner(m.Sender.ID)
		if !fed {
			c.Reply("You aren't the creator of any feds to act in.")
			return nil
		}
		fed_id = fedid
	}
	fed := db.Search_fed_by_id(fed_id)
        if fed["user_id"].(int64) != c.Sender().ID{
          c.Reply("Only the fed creator can export the ban list.")
        return nil
}
	mode := "csv"
	if c.Message().Payload != string("") && stringInSlice(c.Message().Payload, []string{"csv", "json", "xml"}) {
		mode = c.Message().Payload
	}
	if mode == "json" {
           c.Reply("Json" + fmt.Sprint(fed))
	}
	return nil
}




// soon
