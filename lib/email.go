package lib

import "gopkg.in/gomail.v2"

func SendVerifyMail(to string, data string) (err error) {
	m := gomail.NewMessage()
	m.SetHeader("From", "s412172010@gmail.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verify your online teach account!")
	m.SetBody("text/html", `
	<h1>歡迎註冊英文家教平台。</h1>
	<h1>Welcome to register for the OnlineTeach.Asia.</h1>
	<p>在此平台內，您可以自由的選課，找尋想要的外籍英文老師。<br>此平台的目標為降低學生的上課花費以及提升老師的獲利。<br>
	但同時學生以及老師必須自行負擔上課的教材以及其他細項。<br>上課推薦使用Skype語音以保持通訊穩定，以及使用Google Docs來記錄上課資料。
	</p>
	<p>最後如果您同意這些理念，請點擊網址用於開通您的帳號。<a href=https://www.onlineteach.asia/verify/`+data+">https://www.onlineteach.asia/verify/"+data+"</a>"+
	"<p>Please click on the following URL to open your account. <a href=https://www.onlineteach.asia/verify/"+data+">https://www.onlineteach.asia/verify/"+data+"</a></p>")

	d := gomail.NewDialer("smtp.gmail.com", 587, "s412172010@gmail.com", "mlslgurskkgsjdhh")

	err = d.DialAndSend(m)
	return
}
