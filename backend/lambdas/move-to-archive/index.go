package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/glue"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Request struct {
	Detail GlueJobStateChangeEventDetail `json:"detail"`
}

type GlueJobStateChangeEventDetail struct {
	JobName  string `json:"jobName"`
	JobRunID string `json:"jobRunId"`
	State    string `json:"state"`
}

type Response struct {
	S3Key    string `json:"s3Key"`
	S3Bucket string `json:"s3Bucket"`
}

// Event handler, this function handles requests from clients
func HandleInfoEvent(request Request) (Response, error) {
	// Log the event
	log.Println("Received event: ", fmt.Sprintf("%+v", request))

	// Create a new Glue and S3 client to interact with AWS
	log.Println("Creating a new Glue client...")
	mySession := session.Must(session.NewSession())
	glueSvc := glue.New(mySession)
	s3Svc := s3.New(mySession)

	// Get the Glue job status so we can check if it succeeded to determine
	// if we should move the file to the archive bucket or the failed bucket
	log.Println("Getting the Glue job status: ", request.Detail.JobName)
	glueJob, err := glueSvc.GetJobRun(&glue.GetJobRunInput{
		JobName: aws.String(request.Detail.JobName),
		RunId:   aws.String(request.Detail.JobRunID),
	})

	if err != nil {
		log.Println("Error getting Glue job status: ", err)
		return Response{}, err
	}

	// Gets the file processed by the Glue job so we can know what object to copy.
	log.Println("Getting the file processed by the Glue job: ", request.Detail.JobName)
	s3key := glueJob.JobRun.Arguments["--s3_key"]
	s3Bucket := glueJob.JobRun.Arguments["--s3_bucket"]

	// Move the file to the proper folder so the client can know if the file was
	// processed successfully or not
	folder := "archive"
	if request.Detail.State == "FAILED" {
		folder = "failed"
	}

	// The new file key with the new folder.
	newS3Key := fmt.Sprintf("%s/%s", folder, strings.Join(strings.Split(*s3key, "/")[1:], "/"))

	log.Println(fmt.Sprintf("Moving the file to the %s subfolder: ", folder), *s3key)
	_, err = s3Svc.CopyObject(&s3.CopyObjectInput{
		Bucket:     s3Bucket,
		CopySource: aws.String(fmt.Sprintf("%s/%s", *s3Bucket, *s3key)),
		Key:        aws.String(newS3Key),
	})

	if err != nil {
		log.Println("Error moving the file to the archive subfolder: ", err)
		return Response{}, err
	}

	// Delete the file from the original folder
	log.Println("Deleting the file from the original folder: ", request.Detail.JobName)
	_, err = s3Svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: s3Bucket,
		Key:    s3key,
	})

	return Response{S3Key: *s3key, S3Bucket: *s3Bucket}, nil
}

func main() {
	lambda.Start(HandleInfoEvent)
}
