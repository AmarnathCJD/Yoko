package modules


import "gopkg.in/tucnak/telebot.v3"

var ANON = make(map[int]func(tb.Context) error)

func AnonAdmin(next tb.HandlerFunc) tb.HandlerFunc {
fmt.Println("Hmm")
return nil

}

// Soon, today
