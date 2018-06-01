package lib

import "gopkg.in/gomail.v2"

func SendVerifyMail(to string, data string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "s412172010@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify your online teach account!")
	m.SetBody("text/html", "Click this link to verify your account <a href=\"http://localhost:8080/verify/"+data+"\">Verify</a> 。!")

	d := gomail.NewDialer("smtp.gmail.com", 587, "s412172010@gmail.com", "mlslgurskkgsjdhh")

	err = d.DialAndSend(m)
	return
}
