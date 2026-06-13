// Package domain provides the core business models.
package domain

import (
	"fmt"
	"time"
	"unicode/utf8"

	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

// Event represents the core business entity of a event in the system.
// It contains all the essential data and business logic tied to a event.
type Event struct {
	ID      int
	Version int
	Event   string
	Date    time.Time
	UserID  int
}

// NewEvent reconstitues an existing Event entity from storage
// with a known ID and Version.
func NewEvent(
	id int,
	version int,
	event string,
	date time.Time,
	userID int,
) Event {
	return Event{
		ID:      id,
		Version: version,
		Event:   event,
		Date:    date,
		UserID:  userID,
	}
}

// Validate checks whether the business rules for the Event entity are met.
func (e *Event) Validate() error {
	if eventLen := utf8.RuneCountInString(e.Event); eventLen == 0 {
		return fmt.Errorf(
			"invalid 'Event' length %d: %w",
			eventLen,
			errs.ErrInvalidArgument,
		)
	}

	now := time.Now()
	if e.Date.Before(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())) {
		return fmt.Errorf(
			"'Date' can't be before now (%s): %w",
			now.Format("2006-01-02"),
			errs.ErrBadArgument,
		)
	}

	return nil
}

// EventPatch represents the data used to partially update an existing Event.
// If a fiels equal nil, it means it is not applied, in other case it is.
type EventPatch struct {
	Event *string
	Date  *time.Time
}

// NewEventPatch creates a new instance of NewEventPatch.
func NewEventPatch(
	event *string,
	date *time.Time,
) *EventPatch {
	return &EventPatch{
		Event: event,
		Date:  date,
	}
}

// ApplyPatch modifies an Event entity using the provided EventPatch data with business logic validation.
func (e *Event) ApplyPatch(eventPatch EventPatch) error {
	tmp := *e

	if eventPatch.Event != nil {
		tmp.Event = *eventPatch.Event
	}

	if eventPatch.Date != nil {
		tmp.Date = *eventPatch.Date
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf(
			"validate patched event: %w",
			err,
		)
	}

	*e = tmp

	return nil
}
