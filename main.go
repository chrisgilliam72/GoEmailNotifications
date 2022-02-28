package main

import (
	"EmailNotificationService/EmailNotificationLibrary"
	"fmt"
	"log"
)

func main() {
	msgs, err := EmailNotificationLibrary.GetNotifications()

	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < len(msgs); i++ {
		fmt.Printf("Decoded message : %v\n", msgs[i])
	}

}
