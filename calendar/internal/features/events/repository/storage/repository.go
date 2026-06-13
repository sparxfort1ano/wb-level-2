// Package storage acts as the repository layer for the Event feature.
// It uses a build-in map to store events.
package storage

import (
	"sync"
)

type eventID = int

// EventsRepository provides data access methods for task entities.
type EventsRepository struct {
	rw         sync.RWMutex
	storage    map[eventID]EventModel //
	newEventID eventID
	newUserID  int
}

const (
	startEventID = 1
	startUserID  = 0
	startVersion = 1
)

// NewEventsRepository creates a new instance of EventsRepository.
func NewEventsRepository() *EventsRepository {
	return &EventsRepository{
		rw:         sync.RWMutex{},
		storage:    make(map[eventID]EventModel),
		newEventID: startEventID,
		newUserID:  startUserID,
	}
}
