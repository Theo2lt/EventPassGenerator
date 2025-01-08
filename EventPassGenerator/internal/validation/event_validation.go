package validation

import (
	"fmt"
	"time"

	"EventPassGenerator/internal/model"
)

func ConvertRFC3339ToCustom(date string) (string, error) {
	t, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return "", err
	}
	return t.Format("Mon, Jan 02, 2006, 03:04 PM"), nil
}

func BuildValidatedEvent(name, description, location, startDate, endDate string) (*model.Event, error) {
	if len(name) == 0 || len(name) > 50 {
		return nil, fmt.Errorf("name length is invalid")
	}
	if len(description) == 0 || len(description) > 100 {
		return nil, fmt.Errorf("description length is invalid")
	}
	if len(location) == 0 || len(location) > 100 {
		return nil, fmt.Errorf("location length is invalid")
	}

	start, err := time.Parse(time.RFC3339, startDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date")
	}
	end, err := time.Parse(time.RFC3339, endDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date")
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	startFormatted, _ := ConvertRFC3339ToCustom(startDate)
	endFormatted, _ := ConvertRFC3339ToCustom(endDate)

	return &model.Event{
		Name:        name,
		Description: description,
		Location:    location,
		StartDate:   startFormatted,
		EndDate:     endFormatted,
	}, nil
}
