package main

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sqs.New(sess)
	result, err := svc.ListQueues(nil)
	if err != nil {
		logrus.Fatal(err)
	}
	url := *result.QueueUrls[0]
	for {
		output, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: &url,
		})
		if err != nil {
			logrus.Error(err)
		}
		for _, message := range output.Messages {
			logrus.Info(*message.Body)
			_, err = svc.DeleteMessage(&sqs.DeleteMessageInput{
				ReceiptHandle: message.ReceiptHandle,
				QueueUrl:      &url,
			})
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}