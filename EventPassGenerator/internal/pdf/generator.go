package pdf

import (
	"bytes"
	"fmt"
	"image"
	"os"
	"strings"

	"EventPassGenerator/internal/model"
	"EventPassGenerator/internal/validation"

	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
)

func CreatePDF(event *model.Event) ([]byte, error) {
	// 1) Initialize the PDF document (A4 size)
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// 2) Load fonts
	fontFiles := map[string]string{
		"Regular":         "resources/fonts/OpenSans-Regular.ttf",
		"Bold":            "resources/fonts/OpenSans-Bold.ttf",
		"Italic":          "resources/fonts/OpenSans-Italic.ttf",
		"BoldItalic":      "resources/fonts/OpenSans-BoldItalic.ttf",
		"ExtraBold":       "resources/fonts/OpenSans-ExtraBold.ttf",
		"ExtraBoldItalic": "resources/fonts/OpenSans-ExtraBoldItalic.ttf",
		"Light":           "resources/fonts/OpenSans-Light.ttf",
		"LightItalic":     "resources/fonts/OpenSans-LightItalic.ttf",
		"Semibold":        "resources/fonts/OpenSans-Semibold.ttf",
		"SemiboldItalic":  "resources/fonts/OpenSans-SemiboldItalic.ttf",
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

	fmt.Printf("Page width: %.2f, Page height: %.2f\n", pageWidth, pageHeight)

	// 4) For each reservation, create a NEW page
	for i, reservation := range event.Reservations {
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
		imagePath := "resources/images/stockholm.jpg" // Path to the image

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

		fmt.Printf("Image width: %.2f, Image height: %.2f\n", originalWidth, originalHeight)

		ratio := originalWidth / originalHeight

		fmt.Printf("Image ratio: %.2f\n", ratio)
		desiredWidth := totalWidth
		desiredHeight := desiredWidth / ratio

		fmt.Printf("Desired width: %.2f, Desired height: %.2f\n", desiredWidth, desiredHeight)

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
		err = pdf.Cell(nil, strings.ToUpper(reservation.ReservationType))
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

		ReservationAt, err := validation.ConvertRFC3339ToCustom(reservation.ReservationAt)

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
		titleBlock4 := fmt.Sprintf("Ticket %d/%d", i+1, len(event.Reservations))
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
		err = pdf.Cell(nil, fmt.Sprintf("N° %s - Price: %s€", reservation.TicketNumber, reservation.Price))
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
