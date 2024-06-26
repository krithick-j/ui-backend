package service

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"time"
	"ui-back-end/configs"
	"ui-back-end/src/dto"
	"ui-back-end/src/middleware"
	"ui-back-end/src/tmplts"

	"github.com/wneessen/go-mail"
)

var adminmail string = "j.krithick@gmail.com"

func GenOPT() string {
	bi, err := rand.Int(
		rand.Reader,
		big.NewInt(int64(math.Pow(10, 6))),
	)
	if err != nil {
		configs.Log.Errorln("Error on rand Int fn from GenOPT service fn", err.Error())
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
		configs.Log.Errorln("Error on DialAndSend fn from SendPlainMail service fn", err.Error())
		fmt.Println(err.Error())
	}
}

func SendHtmlMailICouopon(tomail string, subject string, data []dto.SendCoupon) {
	//Add admin mail to every email
	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(tomail)
	m.Bcc(adminmail)
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

func SendHtmlMailOrder(order dto.OrderDetailsOut, attachment string) error {

	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(order.CustomerDetails.Email)
	m.Bcc(adminmail)
	m.Subject("Your order is received")
	order.CreatedAt = middleware.FormatTimeByLocation(time.Now(), "Asia/Kolkata", "02-01-2006")

	// funcMap := template.FuncMap{
	// 	"incIndex": incIndex,
	// }
	m.AttachFile(attachment)

	t, err := template.New("email").Parse(tmplts.OrderTemplate)
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
