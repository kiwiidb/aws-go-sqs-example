package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/sirupsen/logrus"
)

type Item struct {
	Message string `json:"message"`
}

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
	result, err := DBsvc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		logrus.Fatal(err)
	}
	if len(result.TableNames) == 0 {
		//create table
		input := &dynamodb.CreateTableInput{
			AttributeDefinitions: []*dynamodb.AttributeDefinition{
				{
					AttributeName: aws.String("message"),
					AttributeType: aws.String("S"),
				},
			},
			KeySchema: []*dynamodb.KeySchemaElement{
				{
					AttributeName: aws.String("message"),
					KeyType:       aws.String("HASH"),
				},
			},
			ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(10),
				WriteCapacityUnits: aws.Int64(10),
			},
			TableName: aws.String(tablename),
		}
		_, err = DBsvc.CreateTable(input)
		if err != nil {
			logrus.Fatal(err)
		}
	}

	url := res.QueueUrl
	for {
		output, err := Queusvc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: url,
		})
		if err != nil {
			logrus.Error(err)
		}
		for _, message := range output.Messages {
			msg := Item{Message: *message.Body}
			av, err := dynamodbattribute.MarshalMap(msg)
			if err != nil {
				logrus.Fatal(err)
			}
			dbInput := dynamodb.PutItemInput{
				Item:      av,
				TableName: &tablename,
			}
			_, err = DBsvc.PutItem(&dbInput)
			if err != nil {
				logrus.Fatal(err)
			}

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
