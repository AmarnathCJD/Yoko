package modules


import "gopkg.in/tucnak/telebot.v3"

var ANON = make(map[int]func(tb.Context) error)

type Update struct {
Id int
Func func(tb.Context) error
Right string


}


func AnonAdmin(next tb.HandlerFunc) tb.HandlerFunc {
fmt.Println("Hmm")
return nil

}

// Soon, today
