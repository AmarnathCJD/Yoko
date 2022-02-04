package modules


import "gopkg.in/tucnak/telebot.v3"

ANON = make(map[int]func(tb.Context) error)

func AnonAdmin(next tb.HandlerFunc) tb.HandlerFunc {
fmt.Println("Hmm")
}

// Soon, today
