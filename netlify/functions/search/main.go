package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mr-destructive/mr-destructive.github.io/models"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func main() {
	lambda.Start(handler)
}

type QueryBody struct {
	Query  string `json:"query"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

func handler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx := context.Background()
	dbName := os.Getenv("TURSO_DATABASE_NAME")
	dbToken := os.Getenv("TURSO_DATABASE_READ_TOKEN")

	var err error
	dbString := fmt.Sprintf("%s?authToken=%s", dbName, dbToken)
	db, err := sql.Open("libsql", dbString)
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Database connection failed"), nil
	}
	defer db.Close()
	body := req.Body
	log.Printf("Body: %s", body)

	var queryBody QueryBody
	err = json.Unmarshal([]byte(body), &queryBody)
	if err != nil {
		return errorResponse(http.StatusBadRequest, "Invalid request body"), nil
	}
	log.Printf("Query: %s", queryBody.Query)

	rows, err := db.QueryContext(ctx, queryBody.Query)
	columns, err := rows.Columns()
	if err != nil {
		return errorResponse(http.StatusInternalServerError, "Database query failed"), nil
	}
	results := map[string]interface{}{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return errorResponse(http.StatusInternalServerError, "Database query failed"), nil
		}
		for i, col := range columns {
			results[col] = values[i]
		}
	}
	return jsonResponse(http.StatusOK, results), nil
}

func jsonResponse(statusCode int, data interface{}) events.APIGatewayProxyResponse {
	body, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "POST, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type, Authorization",
		},
		Body: string(body),
	}
}

func errorResponse(statusCode int, message string) events.APIGatewayProxyResponse {
	return jsonResponse(statusCode, map[string]string{"error": message})
}
