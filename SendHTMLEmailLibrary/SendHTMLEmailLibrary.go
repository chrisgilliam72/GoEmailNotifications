package SendHTMLEmailLibrary

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/gomail.v2"
)

func loadHTMLEmailTemplate() (string, error) {
	strngCurDir, _ := os.Getwd()
	fmt.Printf(" reading html tmplate from %s\n", strngCurDir)
	data, err := ioutil.ReadFile(".\\SendHTMLEmailLibrary\\dealprogressemailtemplate.html")
	if err != nil {
		return "a", fmt.Errorf(" unable to read html template file: %v", err)
	}

	return string(data), nil
}

func SendEmailNotification(emailAddress, emailName, notificationType, applicantionRef, bankReference, eventType string) error {

	emailHtmlTmplate, err := loadHTMLEmailTemplate()
	emailHtmlTmplate = strings.ReplaceAll(emailHtmlTmplate, "{{firstName}}", emailName)
	emailHtmlTmplate = strings.ReplaceAll(emailHtmlTmplate, "{{applicationReference}}", applicantionRef)
	emailHtmlTmplate = strings.ReplaceAll(emailHtmlTmplate, "{{notificationType}}", notificationType)
	emailHtmlTmplate = strings.ReplaceAll(emailHtmlTmplate, "{{bankReference}}", bankReference)
	emailHtmlTmplate = strings.ReplaceAll(emailHtmlTmplate, "{{eventType}}", eventType)
	if err != nil {
		return fmt.Errorf(" unable to load email notification html template:%v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "GoEmailNotificationsManager@comcorp.co.za")
	m.SetHeader("To", emailAddress, emailAddress)
	// m.SetAddressHeader("Cc", "dan@example.com", "Dan")
	m.SetHeader("Subject", "GO Email Notifications")
	m.SetBody("text/html", emailHtmlTmplate)
	// m.Attach("/home/Alex/lolcat.jpg")

	print("Sending email...\n")
	d := gomail.NewDialer("smtp.gmail.com", 587, "chrisgilliam1972@gmail.com", "June1972")
	err = d.DialAndSend(m)
	if err != nil {
		return fmt.Errorf(" unable to send email:%v", err)
	}
	return nil
}
