package main

import (
	"EmailNotificationService/EmailNotificationLibrary"
	"EmailNotificationService/NotificationDatabaseLibrary"
	"EmailNotificationService/SendHTMLEmailLibrary"
	"fmt"
	"log"
	"time"
)

var notificationTypes = map[int]string{
	1: "Deal Progress",
	2: "Feedback Service",
	3: "Template Service",
}

func main() {
	msgCount, err := EmailNotificationLibrary.NotificationCount()
	if err != nil {
		log.Fatalf("Error retrieving msg count %v\n", err)
	}

	fmt.Printf("%d messages waiting \n", msgCount)
	if msgCount > 0 {
		msgs, err := EmailNotificationLibrary.GetNotifications()

		if err != nil {
			log.Fatal(err)
		}

		// for i := 0; i < len(msgs); i++ {
		// 	fmt.Printf("Decoded message : %v\n", msgs[i])
		// }

		for _, msg := range msgs {
			fmt.Printf("Decoded message : %v\n", msg)
			msgId, dbErr := NotificationDatabaseLibrary.AddNotificationMessage(msg.NotificationType, msg.ApplicationReference, msg.BankReference, time.Now(),
				msg.EventType, msg.EventComment, msg.RequestType, msg.MessageStatus, msg.Message, msg.EventId)
			if dbErr != nil {
				log.Fatalf("unable to log message: %v\n", err)
			}

			fmt.Printf(" message %d, logged with database ID :%d/n", msg.EventId, msgId)

			emailAddress, emailName, dbErr := NotificationDatabaseLibrary.GetEmailAddress(msg.ApplicationReference)
			if dbErr != nil {
				fmt.Printf(" unable to retrieve email address:%v", dbErr)
			}

			fmt.Printf("retrieved email address :%s\n", emailAddress)
			fmt.Printf("sending email to :%s\n", emailAddress)
			err = SendHTMLEmailLibrary.SendEmailNotification(emailAddress, emailName, notificationTypes[msg.NotificationType], msg.ApplicationReference, msg.BankReference, msg.EventType)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf(" email sent to  :%s\n", emailAddress)
		}
	}

}
