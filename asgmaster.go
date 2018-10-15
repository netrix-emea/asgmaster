package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// env GOOS=linux  go build -ldflags="-s -w" -o asgmaster
//
// Usage of ./test1:
// -debug
//   	Enable Debugging
// -key string
//   	Tag name to match the asg (default "role")
// -region string
//   	AWS Region (default "eu-west-1")
//
// Usable is cron script   asgmaster && do_action.sh
// will return 0 when there is no other master
//

func main() {

	regionFlag := flag.String("region", "eu-west-1", "AWS Region")
	roleFlag := flag.String("key", "role", "Tag name to match the asg")
	debugFlag := flag.Bool("debug", false, "Enable Debugging")

	flag.Parse()

	if !*debugFlag {
		log.SetOutput(ioutil.Discard)
	}
	// Specify profile for config and region for requests
	awsSession := session.Must(session.NewSessionWithOptions(session.Options{Config: aws.Config{Region: aws.String(*regionFlag)}}))
	ec2svc := ec2.New(awsSession)
	client := ec2metadata.New(awsSession)
	document, _ := client.GetInstanceIdentityDocument()
	if client.Available() {
		//log.Printf(document.InstanceID)
		//log.Printf(document.Region)
	}

	asgRole := getAsgRole(ec2svc, document.InstanceID, *roleFlag)
	masterTag := getMasterTag(ec2svc, asgRole+":master")
	if masterTag == "" {
		// nobody has tag
		log.Printf("No tag available , setting it to myself")
		setMasterTag(ec2svc, document.InstanceID, asgRole+":master")
		os.Exit(0)
	} else if masterTag == document.InstanceID {
		//  we are it
		log.Printf("We are Master")
		os.Exit(0)
	} else {
		// somebody else is it
		log.Printf("Someone else is master")

	}
	log.Printf("Return error")
	os.Exit(1)
}

// getASGRole returns the current ASG role from the current instance
// return "" when not found
func getAsgRole(ec2svc *ec2.EC2, instanceID string, field string) string {
	input := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
			{
				Name: aws.String("resource-type"),
				Values: []*string{
					aws.String("instance"),
				},
			},
			{
				Name: aws.String("key"),
				Values: []*string{
					aws.String(field),
				},
			},
		},
	}

	result, err := ec2svc.DescribeTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Printf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Printf(err.Error())
		}
		return ""
	}

	if len(result.Tags) > 0 {
		return *result.Tags[0].Value
	}
	return ""
}

// returns the MasterTag for all instances with the same ASG Tag
func getMasterTag(ec2svc *ec2.EC2, key string) string {
	input := &ec2.DescribeTagsInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("resource-type"),
				Values: []*string{
					aws.String("instance"),
				},
			},
			{
				Name: aws.String("key"),
				Values: []*string{
					aws.String(key),
				},
			},
		},
	}

	result, err := ec2svc.DescribeTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Printf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Printf(err.Error())
		}
		return ""
	}

	if len(result.Tags) > 0 {
		return *result.Tags[0].Value
	}
	return ""
}

// Creates a tag on the current Instance with name (key) and value (instanceID)
// instanceId = fullname of the isntance
// key = created name normally:   role:master
func setMasterTag(ec2svc *ec2.EC2, instanceID string, key string) string {
	input := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(instanceID),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(key),
				Value: aws.String(instanceID),
			},
		},
	}

	_, err := ec2svc.CreateTags(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				log.Printf(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Printf(err.Error())
		}
		return ""
	}
	//	log.Printf(result)
	return key
}
