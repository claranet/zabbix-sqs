package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {

	type Configuration struct {
		QueueURL        string
		Region          string
		AccessKeyID     string
		SecretAccessKey string
	}

	type Item struct {
		Name  string
		Value string
	}

	type Trigger struct {
		ID       int
		Name     string
		Status   string
		Group    string
		Hostname string
		IP       string
		Severity string
		Items    []Item
	}

	var trigger Trigger
	err := json.Unmarshal([]byte(os.Args[1]), &trigger)
	if err != nil {
		panic(err)
	}

	var items []Item
	for _, item := range trigger.Items {
		if item.Name != "*UNKNOWN*" {
			items = append(items, item)
		}
		fmt.Println(item)
	}
	trigger.Items = items
	messageBytes, err := json.Marshal(trigger)
	if err != nil {
		panic(err)
	}
	message := string(messageBytes[:])
	fmt.Println(message)

	file, _ := os.Open(os.Args[0] + ".json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		panic(err)
	}
	fmt.Println(configuration.QueueURL)

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigDisable,
		//Profile:           "testProfile",
	}))

	svc := sqs.New(sess, &aws.Config{
		Region:      aws.String(configuration.Region),
		Credentials: credentials.NewStaticCredentials(configuration.AccessKeyID, configuration.SecretAccessKey, ""),
	})
	parameters := &sqs.SendMessageInput{
		MessageBody:  aws.String(message),
		QueueUrl:     aws.String(configuration.QueueURL),
		DelaySeconds: aws.Int64(3),
	}
	response, err := svc.SendMessage(parameters)
	if err != nil {
		panic(err)
	}
	fmt.Printf("[Send message] \n%v \n\n", response)

}
