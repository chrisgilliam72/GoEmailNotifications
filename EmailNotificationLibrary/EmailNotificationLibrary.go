package EmailNotificationLibrary

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

type QueueNotificationMessage struct {
	NotificationType     int
	ApplicationReference string
	BankReference        string
	EventDate            string
	EventType            string
	EventComment         string
	RequestType          string
	MessageStatus        string
	Message              string
	EventId              int
}

func (msg QueueNotificationMessage) String() string {
	return fmt.Sprintf(" event ID:%d ApplicationReference:%s BankReference:%s EventDate:%s EventType:%s EventComment:%s RequestType:%s MessageStatus:%s NotificationType:%d",
		msg.EventId, msg.ApplicationReference, msg.BankReference, msg.EventDate, msg.EventType, msg.EventComment, msg.RequestType, msg.MessageStatus, msg.NotificationType)
}

func getAccountURL(storageAccountName, storageQueueName string) (*url.URL, error) {
	_url, err := url.Parse(fmt.Sprintf("https://%s.queue.core.windows.net/%s", storageAccountName, storageQueueName))
	if err != nil {
		err = fmt.Errorf(" error parsing url %v", err)
	}

	return _url, err
}

func getCredentials(storageAccountName, storageAccountKey, storageQueueName string) (azqueue.Credential, error) {

	_url, err := getAccountURL(storageAccountName, storageQueueName)
	if err != nil {
		log.Fatal("Error getting Acount URL: ", err)
	}

	fmt.Printf("Using queue URL: %v\n", _url)

	credential, err := azqueue.NewSharedKeyCredential(storageAccountName, storageAccountKey)
	if err != nil {
		err = fmt.Errorf(" error parsing url %v", err)
	}

	fmt.Printf("Using queue credentials: %v\n", credential.AccountName())

	return credential, err
}

func peekMessages(storageAccountName, storageQueueName string, credential azqueue.Credential) (*azqueue.PeekedMessagesResponse, error) {

	_url, err := getAccountURL(storageAccountName, storageQueueName)
	if err != nil {
		log.Fatal("Error getting Acount URL: ", err)
	}

	fmt.Printf("Peek messages using queue URL: %v\n", _url)

	queueUrl := azqueue.NewQueueURL(*_url, azqueue.NewPipeline(credential, azqueue.PipelineOptions{}))
	ctx := context.TODO()

	props, err := queueUrl.GetProperties(ctx)
	if err != nil {
		// https://godoc.org/github.com/Azure/azure-storage-queue-go/azqueue#StorageErrorCodeType
		errorType := err.(azqueue.StorageError).ServiceCode()

		if errorType == azqueue.ServiceCodeQueueNotFound {

			return nil, fmt.Errorf(" queue not found %v", err)

		} else {
			return nil, fmt.Errorf(" unable to get queue properties %v", err)
		}
	}

	msgUrl := queueUrl.NewMessagesURL()

	messageCount := props.ApproximateMessagesCount()
	if messageCount > 0 {

		// (MessagesURL) Peek(context, maxMessages) (*PeekedMessagesResponse, error)
		peekResp, err := msgUrl.Peek(ctx, 32)
		if err != nil {
			log.Fatal("Error peeking queue messages: ", err)
		}

		log.Printf("Peeked Number of Messages: %d", peekResp.NumMessages())

		return peekResp, nil
	}

	return nil, nil
}

func dequeueMessages(storageAccountName, storageQueueName string, credential azqueue.Credential) (*azqueue.DequeuedMessagesResponse, error) {

	ctx := context.TODO()

	_url, err := getAccountURL(storageAccountName, storageQueueName)
	if err != nil {
		log.Fatal("Error getting Acount URL: ", err)
	}

	fmt.Printf("Dequeue messages using queue URL: %v\n", _url)

	queueUrl := azqueue.NewQueueURL(*_url, azqueue.NewPipeline(credential, azqueue.PipelineOptions{}))
	msgUrl := queueUrl.NewMessagesURL()

	dequeueResp, err := msgUrl.Dequeue(ctx, 32, 10*time.Second)

	if err != nil {
		return nil, fmt.Errorf(" error dequeueing message: %v", err)
	}

	for i := int32(0); i < dequeueResp.NumMessages(); i++ {
		msg := dequeueResp.Message(i)
		log.Printf("Deleting %v: {%v}", i, msg.Text)

		msgIdUrl := msgUrl.NewMessageIDURL(msg.ID)

		// PopReciept is required to delete the Message. If deletion fails using this popreceipt then the message has
		// been dequeued by another client.
		_, err = msgIdUrl.Delete(ctx, msg.PopReceipt)
		if err != nil {
			return nil, fmt.Errorf(" error deleting message: %v", err)
		}
	}
	return dequeueResp, nil
}

func NotificationCount(storageAccountName, storageAccountKey, storageQueueName string) (int, error) {

	credential, err := getCredentials(storageAccountName, storageAccountKey, storageQueueName)
	if err != nil {
		return -1, err
	}

	peekResp, err := peekMessages(storageAccountName, storageQueueName, credential)
	if err != nil {
		return -1, err
	}

	if peekResp == nil {
		return 0, nil
	}
	return int(peekResp.NumMessages()), nil
}

func GetNotifications(storageAccountName, storageAccountKey, storageQueueName string) ([]QueueNotificationMessage, error) {

	credential, err := getCredentials(storageAccountName, storageAccountKey, storageQueueName)
	if err != nil {
		return nil, err
	}

	dequeueMessages, err := dequeueMessages(storageAccountName, storageQueueName, credential)
	if err != nil {
		return nil, err
	}

	var notificationLst = make([]QueueNotificationMessage, dequeueMessages.NumMessages())

	for i := int32(0); i < dequeueMessages.NumMessages(); i++ {
		var queueNotificationMessage QueueNotificationMessage
		msg := dequeueMessages.Message(i)
		json.Unmarshal([]byte(msg.Text), &queueNotificationMessage)
		notificationLst[i] = queueNotificationMessage
	}

	return notificationLst, nil
}
