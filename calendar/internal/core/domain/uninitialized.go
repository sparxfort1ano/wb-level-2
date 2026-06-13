package domain

import "time"

var (
	UninitializedID      = -1
	UninitializedVersion = -1
)

// NewEventUnitialized creates a new Event entity before it is persisted to storage.
// The ID and Version are set to placeholder values.
func NewEventUnitialized(
	event string,
	date time.Time,
	userID int,
) Event {
	return Event{
		ID:      UninitializedID,
		Version: UninitializedVersion,
		Event:   event,
		Date:    date,
		UserID:  userID,
	}
}
