package main

import (
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

// Event handler, this function handles requests from clients
func HandleInfoEvent(request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	ApiResponse := events.APIGatewayV2HTTPResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	log.Println("Received event: ", request)

	// Create a new DynamoDB client
	log.Println("Creating a new DynamoDB client")
	mySession := session.Must(session.NewSession())
	svc := dynamodb.New(mySession)

	// Gets the id from the path
	id := request.PathParameters["id"]

	if id == "" {
		ApiResponse.StatusCode = 400
		body, _ := json.Marshal(&Body{Message: "Error: id is required"})
		ApiResponse.Body = string(body)
		return ApiResponse, nil
	}

	// Convert to float64 to match the type of the Id field in the Item struct
	idFloat, err := strconv.ParseFloat(id, 64)

	// Parse the body of the request into a Item struct
	log.Println("Parsing the body of the request into a Item struct")
	var item Item
	err = json.Unmarshal([]byte(request.Body), &item)
	if err != nil {
		log.Println("Error parsing request body", err)
		ApiResponse.StatusCode = 400
		body, _ := json.Marshal(&Body{Message: "Error: parsing request body"})
		ApiResponse.Body = string(body)
		return ApiResponse, nil
	}

	// If body has an id, check if it matches the id in the path
	if item.Id != 0 && item.Id != idFloat {
		log.Println("Error: id in path does not match id in body")
		ApiResponse.StatusCode = 400
		body, _ := json.Marshal(&Body{Message: "Error: id in path does not match id in body"})
		ApiResponse.Body = string(body)
		return ApiResponse, nil
	} else if item.Id == 0 {
		item.Id = idFloat
	}

	// Convert the Item struct into a DynamoDB AttributeValue map
	log.Println("Converting the Item struct into a DynamoDB AttributeValue map")
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		log.Println("Error marshalling item", err)
		ApiResponse.Body = err.Error()
		ApiResponse.StatusCode = 500
		return ApiResponse, nil
	}

	// Create the DynamoDB PutItemInput object
	log.Println("Creating the DynamoDB PutItemInput object")
	input := &dynamodb.PutItemInput{
		Item:                av,
		TableName:           aws.String(os.Getenv("TABLE_NAME")),
		ConditionExpression: aws.String("attribute_exists(id)"),
	}

	// Write the item to DynamoDB
	log.Println("Writing the item to DynamoDB")
	_, err = svc.PutItem(input)
	if err != nil {
		log.Println("Error writing item to DynamoDB", err)
		body, _ := json.Marshal(&Body{Message: "Error: writing item to DynamoDB"})
		ApiResponse.Body = string(body)
		ApiResponse.StatusCode = 500
		return ApiResponse, nil
	}

	// Return the response to the client
	ApiResponse.StatusCode = 200
	body, _ := json.Marshal(&Body{Message: "Success"})
	ApiResponse.Body = string(body)
	return ApiResponse, nil
}

func main() {
	lambda.Start(HandleInfoEvent)
}
