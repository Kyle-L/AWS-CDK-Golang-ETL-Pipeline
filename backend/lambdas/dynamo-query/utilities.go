package main

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Converts a map of strings to a base64 encoded string
func convertToBase64String(input map[string]*dynamodb.AttributeValue) string {
	jsonString, err := json.Marshal(input)
	if err != nil {
		log.Println("Error converting map to JSON string: ", err)
		return ""
	}
	base64String := base64.StdEncoding.EncodeToString(jsonString)
	log.Println("Base64 encoded string: ", base64String)
	return base64String
}

// Parses a base64 encoded string into a map of strings
func parseBase64String(input string) map[string]*dynamodb.AttributeValue {
	jsonString, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		log.Println("Error decoding base64 string: ", err)
		return nil
	}
	var output map[string]*dynamodb.AttributeValue
	json.Unmarshal(jsonString, &output)
	log.Println("Parsed base64 string: ", output)
	return output
}
