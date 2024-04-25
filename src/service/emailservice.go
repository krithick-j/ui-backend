package service

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"

	"github.com/wneessen/go-mail"
)

func GenOPT() string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, 6))),
	)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%06d", bi)
}

func SendMail(tomail string, subject string, msg string) {
	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(tomail)
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, msg)
	port := mail.WithPort(mail.DefaultPortTLS)
	auth := mail.WithSMTPAuth(mail.SMTPAuthPlain)
	user := mail.WithUsername("api")
	pass := mail.WithPassword("bff432da147b065253a96ea52db78a9d")
	fmt.Println("We are here...")
	fmt.Println("We are here too...")
	go func() {
		client, err := mail.NewClient("live.smtp.mailtrap.io", port, auth, user, pass)
		if err != nil {
			fmt.Println(err.Error())
		}
		defer client.Close()
		client.DialAndSend(m)
	}()
	fmt.Println("We are after email sent")
}
