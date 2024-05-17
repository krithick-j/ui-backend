package service

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"math"
	"math/big"
	"ui-back-end/src/dto"

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
	go func() {
		client := MailFactory()
		defer client.Close()
		client.DialAndSend(m)
	}()
}

func SendHtmlMailICouopon(tomail string, subject string, data []dto.SendCoupon) {
	fmt.Println("Yes, Reached spot 2")
	htmltemplate := `
<!DOCTYPE html>
<html>
<head>
    <title>Your new iCoupon</title>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <h2>This is a noreply email. Your iCoupons are:</h2>
    <table>
        <thead>
            <tr>
                <th>VID</th>
                <th>PIN</th>
                <th>Value</th>
                <th>Date On</th>
                <th>Expires On</th>
                <th>Active</th>
            </tr>
        </thead>
        <tbody>
            {{range .}}
            <tr>
                <td>{{.VID}}</td>
                <td>{{.Pin}}</td>
                <td>{{.Value}}</td>
                <td>{{.DateOn}}</td>
                <td>{{.ExpiresOn}}</td>
                <td>{{.Active}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(tomail)
	m.Subject(subject)
	t, err := template.New("email").Parse(htmltemplate)
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
	htmltemplate := `<!DOCTYPE html>
<html>
<head>
    <title>Your Order Detail</title>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        table {
            width: 100%;
            border-collapse: collapse;
            margin-bottom: 20px;
        }
        th, td {
            border: 1px solid #dddddd;
            text-align: left;
            padding: 8px;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
</head>
<body>
    <p>This is a noreply email. Your Order Details are:</p>
	<p>Dear {{.DeliveryAddress.ContactName}}, </p>
	<p>We received your order on <date>. You have 15 days refund window from the date of order</p>
	<p>Your order details as below</p>
    <table>
        <thead>
            <tr>
                <th>Name</th>
                <th>Qty</th>
                <th>Unit Price</th>
                <th>Total</th>
            </tr>
        </thead>
        <tbody>
            {{range .Items}}
            <tr>
                <td>{{.Name}}</td>
                <td>{{.Quantity}}</td>
                <td>{{.UnitPrice}}</td>
                <td>{{.SubTotal}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`
	fmt.Println("stage 1")
	m := mail.NewMsg()
	m.From("No Reply<admin@ui-network.com>")
	m.To(order.DeliveryAddress.ContactEmail)
	m.Subject("Your order is received")
	t, err := template.New("email").Parse(htmltemplate)
	fmt.Println("stage 2", t)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = m.SetBodyHTMLTemplate(t, order)
	fmt.Println("stage 3")

	if err != nil {
		panic(err.Error())
	}
	client := MailFactory()
	fmt.Println("stage 4")

	defer client.Close()
	err = client.DialAndSend(m)
	fmt.Println("stage 5")

	if err != nil {
		fmt.Println(err.Error())
	}
	return nil
}
