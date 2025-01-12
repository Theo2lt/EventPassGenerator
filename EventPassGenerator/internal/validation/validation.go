package validation

import (
	"EventPassGenerator/internal/model"
)

func ValidatedEvent(event *model.Event) (*model.Event, error) {
	if err := event.Validate(); err != nil {
		return nil, err
	}
	return event, nil
}
