package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

func main() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	Queusvc := sqs.New(sess)
	DBsvc := dynamodb.New(sess)
	tablename := "messages"
	qname := "test-queue"
	res, err := Queusvc.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &qname,
	})
	if err != nil {
		logrus.Fatal(err)
	}

	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tablename),
	}
	_, err = DBsvc.CreateTable(input)
	url := res.QueueUrl
	for {
		output, err := Queusvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: url,
		})
		if err != nil {
			logrus.Error(err)
		}
		for _, message := range output.Messages {
			logrus.Info(*message.Body)
			_, err = Queusvc.DeleteMessage(&sqs.DeleteMessageInput{
				ReceiptHandle: message.ReceiptHandle,
				QueueUrl:      url,
			})
			if err != nil {
				logrus.Error(err)
			}
		}
	}
}
