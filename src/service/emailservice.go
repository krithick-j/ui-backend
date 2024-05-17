package service

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"ui-back-end/src/dto"
	"ui-back-end/src/tmplts"

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

func MailFactory() *mail.Client {
	port := mail.WithPort(mail.DefaultPortTLS)
	auth := mail.WithSMTPAuth(mail.SMTPAuthPlain)
	user := mail.WithUsername("api")
	pass := mail.WithPassword("bff432da147b065253a96ea52db78a9d")
	client, err := mail.NewClient("live.smtp.mailtrap.io", port, auth, user, pass)
	if err != nil {
		fmt.Println(err.Error())
	}
	return client
}

func SendPlainMail(tomail string, subject string, msg string) {
	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(tomail)
	m.Subject(subject)
	m.SetBodyString(mail.TypeTextPlain, msg)
	client := MailFactory()
	defer client.Close()
	err := client.DialAndSend(m)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendHtmlMailICouopon(tomail string, subject string, data []dto.SendCoupon) {

	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(tomail)
	m.Subject(subject)
	t, err := template.New("email").Parse(tmplts.ICouponTemplate)
	if err != nil {
		fmt.Println(err.Error())
	}
	err = m.SetBodyHTMLTemplate(t, data)
	if err != nil {
		panic(err.Error())
	}
	client := MailFactory()
	defer client.Close()
	err = client.DialAndSend(m)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendHtmlMailOrder(order dto.OrderDetailsOut) error {

	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(order.DeliveryAddress.ContactEmail)
	m.Subject("Your order is received")
	t, err := template.New("email").Parse(tmplts.OrderTemplate)
	fmt.Println("stage 2", t)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = m.SetBodyHTMLTemplate(t, order)

	if err != nil {
		panic(err.Error())
	}
	client := MailFactory()

	defer client.Close()
	err = client.DialAndSend(m)

	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
