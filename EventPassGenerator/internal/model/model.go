package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

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

type Event struct {
	Name           string    `json:"name" validate:"required,min=1,max=50"`
	Description    string    `json:"description" validate:"required,min=1,max=60"`
	Location       string    `json:"location" validate:"required,min=1,max=60"`
	StartAt        time.Time `json:"startAt" validate:"required"`
	EndAt          time.Time `json:"endAt" validate:"required,gtfield=StartAt"`
	Reservations   []Person  `json:"reservations" validate:"required,dive"`
	HeaderImageUrl string    `json:"headerImageUrl"`
}

func (e *Event) Validate() error {
	validate := validator.New()
	return validate.Struct(e)
}
