package admission

import (
	"github.com/jung-kurt/gofpdf"
	"gopkg.in/gomail.v2"
)


func MakePDF(){
	gofpdf.Template()
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello, world")
}


func MakeMessage() *gomail.Message{
	m := gomail.NewMessage()
	m.SetHeader("From", "alex@example.com")
	m.SetHeader("To", "bob@example.com", "cora@example.com")
	m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "Hello!")
	m.SetBody("text/html", "Hello <b>Bob</b> and <i>Cora</i>!")
	m.Attach("/home/Alex/lolcat.jpg")
	return m
}


func Run(){
	m := MakeMessage()
	d := gomail.NewDialer("smtp.163.com", 25, "user", "123456")
	// Send the email to Bob, Cora and Dan.
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}