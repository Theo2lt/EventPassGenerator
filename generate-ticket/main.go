package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg" // Required for decoding JPEG images
	_ "image/png"  // Required for decoding PNG images
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

// Person represents an individual with reservation details
type Person struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	ReservationAt   string `json:"reservationAt"` // Example format: "Wed, Mar 12, 2025, 11:00 AM"
	ReservationType string `json:"reservationType"`
	OrderNumber     string `json:"orderNumber"`
	TicketNumber    string `json:"ticketNumber"`
}

// Event represents an event with details and reservations
type Event struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	StartDate   string   `json:"startDate"` // Example format: "Wed, Mar 12, 2025, 11:00 AM"
	EndDate     string   `json:"endDate"`   // Example format: "Wed, Mar 12, 2025, 05:00 PM"
	Reservation []Person `json:"reservation"`
}

// createPDF generates a PDF in memory and returns it as []byte
func createPDF(event Event) ([]byte, error) {
	// 1) Initialize the PDF document (A4 size)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// 2) Load fonts
	fontFiles := map[string]string{
		"Regular":         "fonts/OpenSans-Regular.ttf",
		"Bold":            "fonts/OpenSans-Bold.ttf",
		"Italic":          "fonts/OpenSans-Italic.ttf",
		"BoldItalic":      "fonts/OpenSans-BoldItalic.ttf",
		"ExtraBold":       "fonts/OpenSans-ExtraBold.ttf",
		"ExtraBoldItalic": "fonts/OpenSans-ExtraBoldItalic.ttf",
		"Light":           "fonts/OpenSans-Light.ttf",
		"LightItalic":     "fonts/OpenSans-LightItalic.ttf",
		"Semibold":        "fonts/OpenSans-Semibold.ttf",
		"SemiboldItalic":  "fonts/OpenSans-SemiboldItalic.ttf",
	}

	for name, path := range fontFiles {
		err := pdf.AddTTFFont(name, path)
		if err != nil {
			return nil, fmt.Errorf("failed to add font %s: %v", name, err)
		}
	}

	// 3) Get page dimensions and define margins
	pageWidth, pageHeight := gopdf.PageSizeA4.W, gopdf.PageSizeA4.H
	leftMargin := 0.0
	rightMargin := 0.0
	totalWidth := pageWidth - leftMargin - rightMargin

	// 4) For each reservation, create a NEW page
	for i, reservation := range event.Reservation {
		// --- New page ---
		pdf.AddPage()

		// Default font setting
		err := pdf.SetFont("Regular", "", 12)
		if err != nil {
			return nil, fmt.Errorf("failed to set font: %v", err)
		}

		// --------------------------------------------------------------------
		// (A) Add an image at the top of the page (full width)
		// --------------------------------------------------------------------
		imagePath := "images/event.jpeg" // Path to the image

		// Check if the image file exists
		if _, err := os.Stat(imagePath); os.IsNotExist(err) {
			errMsg := fmt.Sprintf("Image file does not exist: %s", imagePath)
			return nil, fmt.Errorf(errMsg)
		}

		imgFile, err := os.Open(imagePath)
		if err != nil {
			return nil, fmt.Errorf("failed to open image: %v", err)
		}
		defer imgFile.Close()

		// Decode to get dimensions
		imgConfig, _, err := image.DecodeConfig(imgFile)
		if err != nil {
			return nil, fmt.Errorf("failed to decode image: %v", err)
		}
		originalWidth := float64(imgConfig.Width)
		originalHeight := float64(imgConfig.Height)

		ratio := originalWidth / originalHeight
		desiredWidth := totalWidth
		desiredHeight := desiredWidth / ratio

		xPosition := leftMargin
		yPosition := 0.0

		err = pdf.Image(imagePath, xPosition, yPosition, &gopdf.Rect{
			W: desiredWidth,
			H: desiredHeight,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to add image: %v", err)
		}

		// --------------------------------------------------------------------
		// (B) Two side-by-side blocks below the image: block 1 and block 2
		// --------------------------------------------------------------------
		backgroundY := yPosition + desiredHeight
		backgroundHeight := 60.0
		blockWidth := totalWidth / 2

		// ---- BLOCK 1 ----
		pdf.SetFillColor(242, 241, 240)
		pdf.RectFromUpperLeftWithStyle(
			leftMargin,
			backgroundY,
			blockWidth,
			backgroundHeight,
			"F",
		)

		pdf.SetFillColor(0, 0, 0)
		contentY := backgroundY + 10
		leftMarginText := leftMargin + 25.0

		// Title (bold)
		err = pdf.SetFont("Bold", "", 10)
		if err != nil {
			return nil, fmt.Errorf("failed to set Bold font: %v", err)
		}
		pdf.SetX(leftMarginText)
		pdf.SetY(contentY)
		err = pdf.Cell(nil, event.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to add title: %v", err)
		}

		// Subtitle
		contentY += 18
		pdf.SetFont("Regular", "", 9)
		pdf.SetX(leftMarginText)
		pdf.SetY(contentY)
		err = pdf.Cell(nil, event.Description)
		if err != nil {
			return nil, fmt.Errorf("failed to add subtitle: %v", err)
		}

		// ---- BLOCK 2 ----
		pdf.SetFillColor(242, 241, 240)
		pdf.RectFromUpperLeftWithStyle(
			leftMargin+blockWidth,
			backgroundY,
			blockWidth,
			backgroundHeight,
			"F",
		)
		pdf.SetFillColor(0, 0, 0)

		contentY2 := backgroundY + 9
		leftMarginText2 := (leftMargin + blockWidth)

		err = pdf.SetFont("Semibold", "", 7)
		if err != nil {
			return nil, fmt.Errorf("failed to set Semibold font: %v", err)
		}
		pdf.SetX(leftMarginText2)
		pdf.SetY(contentY2)
		err = pdf.Cell(nil, "TICKET VALIDITY")
		if err != nil {
			return nil, fmt.Errorf("failed to add TICKET VALIDITY text: %v", err)
		}

		contentY2 += 10
		err = pdf.SetFont("Semibold", "", 8)
		if err != nil {
			return nil, fmt.Errorf("failed to set Semibold font: %v", err)
		}
		pdf.SetX(leftMarginText2)
		pdf.SetY(contentY2)
		err = pdf.Cell(nil, event.StartDate+" – "+event.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to add date/time text: %v", err)
		}

		contentY2 += 15
		err = pdf.SetFont("Semibold", "", 7)
		if err != nil {
			return nil, fmt.Errorf("failed to set Semibold font: %v", err)
		}
		pdf.SetX(leftMarginText2)
		pdf.SetY(contentY2)
		err = pdf.Cell(nil, "ADDRESS")
		if err != nil {
			return nil, fmt.Errorf("failed to add address text: %v", err)
		}

		contentY2 += 10
		err = pdf.SetFont("Semibold", "", 8)
		if err != nil {
			return nil, fmt.Errorf("failed to set Semibold font: %v", err)
		}
		pdf.SetX(leftMarginText2)
		pdf.SetY(contentY2)
		err = pdf.Cell(nil, event.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to add address: %v", err)
		}

		// --------------------------------------------------------------------
		// (C) Two side-by-side blocks under block 1 and block 2: block 3 & 4
		// --------------------------------------------------------------------
		newRowY := backgroundY + backgroundHeight

		// ---- BLOCK 3 ----
		pdf.SetFillColor(255, 255, 255)
		pdf.RectFromUpperLeftWithStyle(
			leftMargin,
			newRowY,
			blockWidth,
			pageHeight,
			"F",
		)
		pdf.SetFillColor(0, 0, 0)
		contentY3 := newRowY + 20
		leftMarginText3 := leftMargin + 25.0

		err = pdf.SetFont("Bold", "", 10)
		if err != nil {
			return nil, fmt.Errorf("failed to set Bold font: %v", err)
		}
		pdf.SetX(leftMarginText3)
		pdf.SetY(contentY3)
		err = pdf.Cell(nil, "EARLY TICKET")
		if err != nil {
			return nil, fmt.Errorf("failed to add text in block 3: %v", err)
		}
		contentY3 += 25

		pdf.SetFont("Semibold", "", 8)
		pdf.SetX(leftMarginText3)
		pdf.SetY(contentY3)
		err = pdf.Cell(nil, "Event date: "+event.StartDate+" – "+event.EndDate)
		if err != nil {
			return nil, fmt.Errorf("failed to add content in block 3: %v", err)
		}
		contentY3 += 10

		pdf.SetFont("Regular", "", 8)
		pdf.SetX(leftMarginText3)
		pdf.SetY(contentY3)

		ReservationAt, err := ConvertRFC3339ToCustom(reservation.ReservationAt)

		err = pdf.Cell(nil, fmt.Sprintf("Order - N° %s - %s",
			reservation.OrderNumber,
			ReservationAt,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to add content in block 3: %v", err)
		}
		contentY3 += 10

		pdf.SetFont("Regular", "", 8)
		pdf.SetX(leftMarginText3)
		pdf.SetY(contentY3)
		err = pdf.Cell(nil, fmt.Sprintf("Customer: %s %s - %s",
			reservation.FirstName,
			reservation.LastName,
			reservation.Email,
		))
		if err != nil {
			return nil, fmt.Errorf("failed to add content in block 3: %v", err)
		}
		contentY3 += 10

		// ---- BLOCK 4 ----
		pdf.SetFillColor(250, 250, 250)
		pdf.RectFromUpperLeftWithStyle(
			leftMargin+blockWidth,
			newRowY,
			blockWidth,
			pageHeight,
			"F",
		)
		pdf.SetFillColor(0, 0, 0)

		block4Left := leftMargin + blockWidth
		block4Width := blockWidth
		centerBlock4X := block4Left + (block4Width / 2)

		contentY4 := newRowY + 20

		err = pdf.SetFont("Bold", "", 10)
		if err != nil {
			return nil, fmt.Errorf("failed to set Bold font: %v", err)
		}

		// Example: "Ticket 1/5"
		titleBlock4 := fmt.Sprintf("Ticket %d/%d", i+1, len(event.Reservation))
		titleWidth, err := pdf.MeasureTextWidth(titleBlock4)
		if err != nil {
			return nil, fmt.Errorf("failed to measure text width in block 4: %v", err)
		}
		pdf.SetX(centerBlock4X - (titleWidth / 2))
		pdf.SetY(contentY4)
		err = pdf.Cell(nil, titleBlock4)
		if err != nil {
			return nil, fmt.Errorf("failed to add text in block 4: %v", err)
		}

		// --- Generate a QR code for each reservation ---
		contentY4 += 15
		err = pdf.SetFont("Regular", "", 9)
		if err != nil {
			return nil, fmt.Errorf("failed to set Regular font: %v", err)
		}

		qrData := fmt.Sprintf("https://www.linkedin.com/in/theo-liot/")
		qrPNG, err := qrcode.Encode(qrData, qrcode.Medium, 256)
		if err != nil {
			return nil, fmt.Errorf("error generating QR code: %v", err)
		}
		imgHolder, err := gopdf.ImageHolderByBytes(qrPNG)
		if err != nil {
			return nil, fmt.Errorf("error creating image holder: %v", err)
		}
		qrWidth := 140.0
		qrHeight := 140.0
		qrPosX := centerBlock4X - (qrWidth / 2)

		err = pdf.ImageByHolder(imgHolder, qrPosX, contentY4, &gopdf.Rect{W: qrWidth, H: qrHeight})
		if err != nil {
			return nil, fmt.Errorf("error inserting QR code: %v", err)
		}

		// Additional info below the QR code
		pdf.SetFont("Semibold", "", 8)
		pdf.SetX(qrPosX)
		pdf.SetY(contentY4 + qrHeight + 10)
		err = pdf.Cell(nil, fmt.Sprintf("N° %s - Price: 10.00€", reservation.TicketNumber))
		if err != nil {
			return nil, fmt.Errorf("failed to add content in block 4: %v", err)
		}

		// --------------------------------------------------------------------
		// (D) Footer centered at the bottom of the page
		// --------------------------------------------------------------------
		err = pdf.SetFont("Italic", "", 9)
		if err != nil {
			return nil, fmt.Errorf("failed to set Italic font: %v", err)
		}
		footerText := "Note: This ticket is not real!"

		textWidth, err := pdf.MeasureTextWidth(footerText)
		if err != nil {
			return nil, fmt.Errorf("failed to measure text width: %v", err)
		}
		pdf.SetX((pageWidth - textWidth) / 2)
		pdf.SetY(pageHeight - 25)
		err = pdf.Cell(nil, footerText)
		if err != nil {
			return nil, fmt.Errorf("failed to add footer: %v", err)
		}

		// End of the page (loop continues for the next reservation)
	}

	// ------------------------------------------------------------------------
	// (E) Final PDF generation in memory (byte array)
	// ------------------------------------------------------------------------
	var buf bytes.Buffer
	err := pdf.Write(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %v", err)
	}

	return buf.Bytes(), nil
}

// ConvertRFC3339ToCustom converts an RFC3339 string to a custom format "Mon, Jan 02, 2006, 03:04 PM"
func ConvertRFC3339ToCustom(rfc3339Str string) (string, error) {
	t, err := time.Parse(time.RFC3339, rfc3339Str)
	if err != nil {
		return "", err
	}
	desiredFormat := "Mon, Jan 02, 2006, 03:04 PM"
	return t.Format(desiredFormat), nil
}

// NewEvent creates a new Event with validated and formatted fields
func NewEvent(name, description, location string, startDateStr, endDateStr string) (*Event, error) {
	// Define maximum lengths
	const (
		maxNameLength        = 42
		maxDescriptionLength = 54
		maxLocationLength    = 64
	)

	// Validate name length
	if len(name) == 0 {
		return nil, fmt.Errorf("Name cannot be empty")
	}
	if len(name) > maxNameLength {
		return nil, fmt.Errorf("Name exceeds maximum length of %d characters", maxNameLength)
	}

	// Validate description length
	if len(description) == 0 {
		return nil, fmt.Errorf("Description cannot be empty")
	}
	if len(description) > maxDescriptionLength {
		return nil, fmt.Errorf("Description exceeds maximum length of %d characters", maxDescriptionLength)
	}

	// Validate location length
	if len(location) == 0 {
		return nil, fmt.Errorf("Location cannot be empty")
	}
	if len(location) > maxLocationLength {
		return nil, fmt.Errorf("Location exceeds maximum length of %d characters", maxLocationLength)
	}

	// Validate RFC3339 format for start date
	startDate, err := time.Parse(time.RFC3339, startDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid StartDate: %v", err)
	}

	// Validate RFC3339 format for end date
	endDate, err := time.Parse(time.RFC3339, endDateStr)
	if err != nil {
		return nil, fmt.Errorf("invalid EndDate: %v", err)
	}

	// Ensure end date is after start date
	if endDate.Before(startDate) {
		return nil, fmt.Errorf("EndDate must be after StartDate")
	}

	// Convert dates to custom format
	startDateFormatted, err := ConvertRFC3339ToCustom(startDateStr)
	if err != nil {
		return nil, fmt.Errorf("error converting StartDate: %v", err)
	}

	endDateFormatted, err := ConvertRFC3339ToCustom(endDateStr)
	if err != nil {
		return nil, fmt.Errorf("error converting EndDate: %v", err)
	}

	return &Event{
		Name:        name,
		Description: description,
		Location:    location,
		StartDate:   startDateFormatted,
		EndDate:     endDateFormatted,
		Reservation: []Person{}, // Initialize with an empty slice
	}, nil
}

// EventPayload corresponds to the expected JSON input fields
type EventPayload struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Location    string   `json:"location"`
	StartDate   string   `json:"startDate"` // Expecting RFC3339
	EndDate     string   `json:"endDate"`   // Expecting RFC3339
	Reservation []Person `json:"reservation"`
}

// handler is the AWS Lambda (APIGateway) handler function
func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// 1) Parse the JSON request from request.Body
	var payload EventPayload
	if err := json.Unmarshal([]byte(request.Body), &payload); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "The request body is invalid or missing",
		}, nil
	}

	// 2) Create the Event with validation and date formatting
	event, err := NewEvent(
		payload.Name,
		payload.Description,
		payload.Location,
		payload.StartDate,
		payload.EndDate,
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid parameters: %v", err),
		}, nil
	}

	// Add the list of reservations
	event.Reservation = payload.Reservation

	// 3) Generate the PDF
	pdfBytes, err := createPDF(*event)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Internal server error",
		}, nil
	}

	// 4) Base64 encode the PDF
	encoded := base64.StdEncoding.EncodeToString(pdfBytes)

	// 5) Return the HTTP response with the base64-encoded PDF
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Content-Type":        "application/pdf",
			"Content-Disposition": "attachment; filename=document.pdf",
		},
		Body:            encoded,
		IsBase64Encoded: true,
	}, nil
}

func main() {
	// Start the Lambda function
	lambda.Start(handler)
}
