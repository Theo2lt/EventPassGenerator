package model

type Person struct {
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	ReservationAt   string `json:"reservationAt"`
	ReservationType string `json:"reservationType"`
	OrderNumber     string `json:"orderNumber"`
	TicketNumber    string `json:"ticketNumber"`
	Price           string `json:"price"`
}

type Event struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Location     string   `json:"location"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	Reservations []Person `json:"reservations"`
}
