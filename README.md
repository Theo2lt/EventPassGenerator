# EventPassGenerator

## Overview

**EventPassGenerator** is a serverless application built on AWS Lambda that generates event passes in PDF format. The application accepts event and reservation data via an API Gateway request, validates the input according to predefined rules, and produces a personalized PDF that includes event details, reservation details, and an integrated unique QR code. The resulting PDF is returned as a Base64-encoded response.

---

## Features

- **Dynamic PDF Generation**: Produces a personalized event pass in PDF format.
- **Input Validation**: Ensures that all necessary data is provided and meets specific validation rules.
- **QR Code Integration**: Each reservation gets a unique QR code for quick check-ins or validation.
- **Customizable Header Image**: Incorporates a header image into the PDF (with fallback if not provided).
- **Serverless Design**: Utilizes AWS Lambda and API Gateway for scalable and efficient operation.

---

## Prerequisites

1. **Go**: Install Go version 1.23.4.
2. **Terraform**: Ensure you have Terraform installed ([download here](https://www.terraform.io/downloads.html)).
3. **Dependencies**:  
   - [`gopdf`](https://github.com/signintech/gopdf): For PDF creation.
   - [`go-qrcode`](https://github.com/skip2/go-qrcode): For generating QR codes.
   - [`aws-lambda-go`](https://github.com/aws/aws-lambda-go): For integration with AWS Lambda.
   - [`validator/v10`](https://github.com/go-playground/validator): For input validation.

Install all dependencies using:

```bash
go mod tidy
```

---

## Data Structures and Validation Requirements

To ensure proper data handling, the following Go structures and validation tags must be used:

### Person Structure

Represents an individual reservation. Each field must adhere to the following validations:

```go
type Person struct {
	FirstName       string    `json:"firstName" validate:"required"`
	LastName        string    `json:"lastName" validate:"required"`
	Email           string    `json:"email" validate:"required,email"`
	ReservationAt   time.Time `json:"reservationAt" validate:"required"`
	ReservationType string    `json:"reservationType" validate:"required"`
	OrderNumber     string    `json:"orderNumber" validate:"required"`
	TicketNumber    string    `json:"ticketNumber" validate:"required,len=12"`
	Price           string    `json:"price" validate:"required,number"`
}
```

- **FirstName** and **LastName**: Required.
- **Email**: Must be provided and be a valid email address.
- **ReservationAt**: Required date and time of reservation.
- **ReservationType**: Required field to specify the type of reservation.
- **OrderNumber**: Required order identifier.
- **TicketNumber**: Required ticket identifier (exactly 12 characters).
- **Price**: Must be provided and represent a numeric value.

### Event Structure

Represents the event details along with the associated reservations:

```go
type Event struct {
	Name           string    `json:"name" validate:"required,min=1,max=50"`
	Description    string    `json:"description" validate:"required,min=1,max=60"`
	Location       string    `json:"location" validate:"required,min=1,max=60"`
	StartAt        time.Time `json:"startAt" validate:"required"`
	EndAt          time.Time `json:"endAt" validate:"required,gtfield=StartAt"`
	Reservations   []Person  `json:"reservations" validate:"required,dive"`
	HeaderImageUrl string    `json:"headerImageUrl"`
}
```

- **Name**: Required event name (between 1 and 50 characters).
- **Description**: Required description (between 1 and 60 characters).
- **Location**: Required event location (between 1 and 60 characters).
- **StartAt**: Required starting date and time.
- **EndAt**: Required ending date and time, which must be later than `StartAt`.
- **Reservations**: An array of `Person` entries, all of which need to pass validation.
- **HeaderImageUrl**: Optional URL for the header image in the PDF. If omitted, a default image will be used.

---

## Example Payload

Below is an example JSON payload that meets the required structure:

```json
{
  "name": "Tech Conference 2025",
  "description": "A cutting-edge tech conference",
  "location": "Downtown Convention Center",
  "startAt": "2025-06-15T09:00:00Z",
  "endAt": "2025-06-15T17:00:00Z",
  "reservations": [
    {
      "firstName": "John",
      "lastName": "Doe",
      "email": "john.doe@example.com",
      "reservationAt": "2025-06-01T12:00:00Z",
      "reservationType": "VIP",
      "orderNumber": "ORD123456789",
      "ticketNumber": "TICK1234567890",
      "price": "100"
    }
  ],
  "headerImageUrl": "https://example.com/image.jpg"
}
```

---

## PDF Generation and QR Code Integration

1. **Validation**:
   - The payload is validated using the `validator/v10` library.
   - Any missing or invalid fields will result in a `400 Bad Request` error.

2. **PDF Creation**:
   - The PDF document is generated using the `gopdf` library.
   - If a `headerImageUrl` is provided, the image is loaded; otherwise, a default header image is used.
   - A unique QR code is generated for each reservation using the `go-qrcode` library.

3. **Response**:
   - The generated PDF is returned as a Base64-encoded string.
   - Response headers include `Content-Type: application/pdf` and `Content-Disposition: attachment` to ensure proper handling by clients.

---

## Deployment Using Terraform

Instead of deploying with the AWS CLI, deploy the function using Terraform. Follow these steps:

1. Change into the Terraform directory:

   ```bash
   cd terraforme
   ```

2. Initialize Terraform (if you haven’t already):

   ```bash
   terraform init
   ```

3. Apply the Terraform configuration:

   ```bash
   terraform apply
   ```

Terraform will handle the creation of the Lambda function, API Gateway, IAM roles, and any other resources defined in your Terraform configuration files.

---

## Configuration

### Environment Variables

Configure your Lambda function with the following environment variables:

- **Default Image Path**: Specify the fallback header image to be used if no `headerImageUrl` is provided.
- **Resource Paths**: Ensure the correct paths to resources such as `resources/images` and `resources/fonts` are set up within your deployment package.

---

## Image Requirements

The header image provided should adhere to the following specifications:

### Aspect Ratio

- **1.78 (16:9)**: The image's width should be 1.78 times larger than its height.

### Example Resolutions

- **1280 x 720 (HD)**
- **1920 x 1080 (Full HD)**
- **2560 x 1440 (Quad HD)**
- **3840 x 2160 (4K UHD)**
- **7680 x 4320 (8K UHD)**

### Format

- The image must be in **JPEG** or **JPG** format.

---

## Testing the API

You can test the **EventPassGenerator** using the following test endpoint:

**Test Endpoint:**  
[https://ekqcygxc19.execute-api.eu-west-1.amazonaws.com/](https://ekqcygxc19.execute-api.eu-west-1.amazonaws.com/)

### Testing Instructions

1. **Method**: Use the `POST` method to send requests.

2. **cURL Example**:

   ```bash
   curl -X POST \
     -H "Content-Type: application/json" \
     -d '{
       "name": "Tech Conference 2025",
       "description": "A cutting-edge tech conference",
       "location": "Downtown Convention Center",
       "startAt": "2025-06-15T09:00:00Z",
       "endAt": "2025-06-15T17:00:00Z",
       "reservations": [
         {
           "firstName": "John",
           "lastName": "Doe",
           "email": "john.doe@example.com",
           "reservationAt": "2025-06-01T12:00:00Z",
           "reservationType": "VIP",
           "orderNumber": "ORD123456789",
           "ticketNumber": "TICK1234567890",
           "price": "100"
         }
       ],
       "headerImageUrl": "https://example.com/image.jpg"
     }' \
     https://ekqcygxc19.execute-api.eu-west-1.amazonaws.com/
   ```

3. **Postman Instructions**:
   - Open Postman.
   - Set the request type to `POST`.
   - Enter the test endpoint: `https://ekqcygxc19.execute-api.eu-west-1.amazonaws.com/`
   - In the **Headers** tab, set `Content-Type` to `application/json`.
   - In the **Body** tab, select **raw** and paste the JSON payload.
   - Click **Send**.

4. **Response**:  
   The API will return a Base64-encoded PDF. Clients should decode the Base64 string to view or store the generated PDF. The response includes appropriate headers such as:
   - `Content-Type: application/pdf`
   - `Content-Disposition: attachment`

---

## Error Handling

- **400 Bad Request**: Returned when the incoming request payload is invalid or missing required fields.
- **500 Internal Server Error**: Returned when an error occurs during PDF generation or image processing.

---

## Libraries Used

- [gopdf](https://github.com/signintech/gopdf) – PDF generation library.
- [go-qrcode](https://github.com/skip2/go-qrcode) – QR code generation library.
- [aws-lambda-go](https://github.com/aws/aws-lambda-go) – AWS Lambda Go SDK.
- [validator/v10](https://github.com/go-playground/validator) – Input validation library.

---