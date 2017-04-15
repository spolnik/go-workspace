package main

import (
	"fmt"
)

type user struct {
	name  string
	email string
}

func (u user) notify() {
	fmt.Printf("Sending User Email To %s<%s>\n",
		u.name,
		u.email)
}

func (u *user) changeEmail(email string) {
	u.email = email
}

func (u user) valueChangeEmail(email string) {
	u.email = email
}

func main() {
	bill := user{"Bill", "bill@email.com"}
	bill.notify()

	lisa := &user{"Lisa", "lisa@email.com"}
	lisa.notify()

	bill.changeEmail("bill@newdomain.com")
	bill.notify()

	lisa.changeEmail("lisa@newdomain.com")
	lisa.notify()

	// as below is on value receiver, done on copy - it will not work on original receiver

	bill.valueChangeEmail("bill@totallyNewDomain.com")
	bill.notify()

	lisa.valueChangeEmail("lisa@totallyNewDomain.com")
	lisa.notify()
}
