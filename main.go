package main

import (
	"EmailNotificationService/EmailNotificationLibrary"
	"EmailNotificationService/NotificationDatabaseLibrary"
	"EmailNotificationService/SendHTMLEmailLibrary"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

var notificationTypes = map[int]string{
	1: "Deal Progress",
	2: "Feedback Service",
	3: "Template Service",
}

type Config struct {
	QueueSettings struct {
		StorageAccountName string `json:"StorageAccountName"`
		StorageAccountKey  string `json:"StorageAccountKey"`
		StorageQueueName   string `json:"StorageQueueName"`
	} `json:"QueueSettings"`
}

func (cfg Config) String() string {
	return fmt.Sprintf("Account Name: %s Account Key: %s QueueName: %s", cfg.QueueSettings.StorageAccountName, cfg.QueueSettings.StorageAccountKey, cfg.QueueSettings.StorageQueueName)
}

func loadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return Config{}, fmt.Errorf(" unable to open config file: %v", err)
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config, nil
}

func main() {

	var config Config
	var err error
	config, err = loadConfiguration("config.json")
	if err != nil {
		log.Fatalf("Error retrieving application settings  %v\n", err)
	}
	msgCount, err := EmailNotificationLibrary.NotificationCount(config.QueueSettings.StorageAccountName, config.QueueSettings.StorageAccountKey, config.QueueSettings.StorageQueueName)
	if err != nil {
		log.Fatalf("Error retrieving msg count %v\n", err)
	}

	fmt.Printf("%d messages waiting \n", msgCount)
	if msgCount > 0 {
		msgs, err := EmailNotificationLibrary.GetNotifications(config.QueueSettings.StorageAccountName, config.QueueSettings.StorageAccountKey, config.QueueSettings.StorageQueueName)

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
