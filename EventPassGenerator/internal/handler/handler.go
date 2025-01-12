package handler

import (
	"EventPassGenerator/internal/model"
	"EventPassGenerator/internal/pdf"
	"EventPassGenerator/internal/validation"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
)

func LambdaHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var payload model.Event
	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid or missing request body",
		}, nil
	}

	event, err := validation.BuildValidatedEvent(
		payload.Name, payload.Description, payload.Location, payload.StartDate, payload.EndDate,
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Validation error: %v", err),
		}, nil
	}

	event.Reservations = payload.Reservations

	pdfBytes, err := pdf.CreatePDF(event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to generate PDF",
		}, nil
	}

	encoded := base64.StdEncoding.EncodeToString(pdfBytes)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":        "application/pdf",
			"Content-Disposition": "attachment; filename=event_pass.pdf",
		},
		Body:            encoded,
		IsBase64Encoded: true,
	}, nil
}
