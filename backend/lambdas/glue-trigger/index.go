package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
)

type Request struct {
	Detail GlueTriggerEventDetail `json:"detail"`
}

type GlueTriggerEventDetail struct {
	Bucket S3Bucket `json:"bucket"`
	Object S3Object `json:"object"`
}

type S3Bucket struct {
	Name string `json:"name"`
}

type S3Object struct {
	Key string `json:"key"`
}

// Event handler, this function handles requests from clients
func HandleInfoEvent(request Request) (string, error) {
	// Log the event
	log.Println("Received event: ", fmt.Sprintf("%+v", request))

	// Check to see if the S3 key prefix is correct
	if !strings.HasPrefix(request.Detail.Object.Key, os.Getenv("S3_KEY_PREFIX")) {
		log.Println(fmt.Sprintf("S3 key prefix does not match the desired prefix: %s", os.Getenv("S3_KEY_PREFIX")))
		return "", nil
	}

	// Check the file is a .csv file, we don't want to process any other files
	if !strings.HasSuffix(request.Detail.Object.Key, ".csv") {
		log.Println(fmt.Sprintf("S3 key does not match the desired file type: %s", ".csv"))
		return "", nil
	}

	// Create a new Glue client
	log.Println("Creating a new Glue client...")
	mySession := session.Must(session.NewSession())
	svc := glue.New(mySession)

	// Create a new Glue job
	log.Println("Creating a new Glue job: ", os.Getenv("JOB_NAME"))
	glueJob, err := svc.StartJobRun(&glue.StartJobRunInput{
		JobName: aws.String(os.Getenv("JOB_NAME")),
		Arguments: map[string]*string{
			"--s3_bucket": aws.String(request.Detail.Bucket.Name),
			"--s3_key":    aws.String(request.Detail.Object.Key),
			"--table":     aws.String(os.Getenv("TABLE_NAME")),
			"--workers":   aws.String(os.Getenv("WORKERS")),
		},
	})

	if err != nil {
		log.Println("Error starting Glue job: ", err)
		return "", err
	}

	// Return the response to the client
	return aws.StringValue(glueJob.JobRunId), nil
}

func main() {
	lambda.Start(HandleInfoEvent)
}
