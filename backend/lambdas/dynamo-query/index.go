package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	maxPageSize = 250
)

// Event handler, this function handles requests from clients
func HandleInfoEvent(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ApiResponse := events.APIGatewayV2HTTPResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	log.Println("Received event: ", request)

	// Query parameters
	day := request.QueryStringParameters["day"]
	if day != "" {
		dayInt, _ := strconv.Atoi(day)
		day = fmt.Sprintf("%02d", dayInt)
	}
	month, _ := strconv.Atoi(request.QueryStringParameters["month"])
	monthString := fmt.Sprintf("%02d", month)
	year := request.QueryStringParameters["year"]
	pageSize := request.QueryStringParameters["pageSize"]

	// Set default page size
	pageSizeInt := int64(maxPageSize)
	if pageSize != "" {
		pageSizeInt, _ = strconv.ParseInt(pageSize, 10, 64)
		if pageSizeInt > maxPageSize {
			pageSizeInt = maxPageSize
		}
	}

	// If Fraud is present, convert to string so it can be used with the
	// global secondary index. If not set, defaults to "FALSE".
	isFraud := request.QueryStringParameters["isFraud"] == "true"
	isFraudString := "FALSE"
	if isFraud {
		isFraudString = "TRUE"
	}

	// Get pagination last evaluated key & limit from query string
	paginationToken := parseBase64String(request.QueryStringParameters["paginationToken"])
	limit, _ := strconv.ParseInt(request.QueryStringParameters["limit"], 10, 64)

	// Set default limit to 100
	if limit == 0 || limit > maxPageSize {
		limit = maxPageSize
	}

	// Create a new DynamoDB client
	log.Println("Creating a new DynamoDB client")
	mySession := session.Must(session.NewSession())
	svc := dynamodb.New(mySession)

	// Create a new DynamoDB query using global secondary index
	log.Println("Creating a new DynamoDB query")
	dynamoQuery, err := svc.Query(&dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("TABLE_NAME")),
		IndexName:              aws.String(os.Getenv("INDEX_NAME")),
		ConsistentRead:         aws.Bool(false),
		KeyConditionExpression: aws.String("#isFraud = :isFraud AND begins_with(#transactionDateTime, :transactionDateTime)"),
		ExpressionAttributeNames: map[string]*string{
			"#isFraud":             aws.String("isFraud"),
			"#transactionDateTime": aws.String("transactionDateTime"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":isFraud":             {S: aws.String(isFraudString)},
			":transactionDateTime": {S: aws.String(fmt.Sprintf("%s-%s-%s", year, monthString, day))},
		},
		Limit:             aws.Int64(pageSizeInt),
		ExclusiveStartKey: paginationToken,
	})

	if err != nil {
		log.Println("Error querying DynamoDB: ", err)
		body := fmt.Sprintf("Error querying DynamoDB: %s", err)
		ApiResponse.Body = body
		ApiResponse.StatusCode = 500
		return ApiResponse, nil
	}

	// Unmarshall the response into a slice of Item structs
	// This removes the types from the DynamoDB response
	log.Println("Unmarshalling the response into a slice of Item structs")
	var items []Item
	err = dynamodbattribute.UnmarshalListOfMaps(dynamoQuery.Items, &items)

	if err != nil {
		log.Println("Error formatting DynamoDB response: ", err)
		body := fmt.Sprintf("Error formatting DynamoDB response: %s", err)
		ApiResponse.Body = body
		ApiResponse.StatusCode = 500
		return ApiResponse, nil
	}

	// Converts the lastKeyEvaluated into a base64 encoded string
	// This is used for pagination
	log.Println("Converting the lastKeyEvaluated into a base64 encoded string")
	lastEvaluatedKeyString := ""
	if dynamoQuery.LastEvaluatedKey != nil {
		lastEvaluatedKeyString = convertToBase64String(dynamoQuery.LastEvaluatedKey)
	}

	// Marshall the slice of Item structs into a JSON string
	// This adds the field names back into the response and makes it easier to read
	// for the client.
	log.Println("Marshalling the slice of Item structs into a JSON string")
	body := &map[string]interface{}{
		"items":           items,
		"count":           len(items),
		"paginationToken": lastEvaluatedKeyString,
	}
	json, err := json.Marshal(body)

	// Return the response to the client
	ApiResponse.Body = string(json)
	ApiResponse.StatusCode = 200
	return ApiResponse, nil
}

func main() {
	lambda.Start(HandleInfoEvent)
}
