// Package service acts as the service layer for the Event feature.
// Its is responsible for validating the event payload.
package service

import (
	"context"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

// EventsService encapsulates the core business logic for event management.
// All of its methods delegate the persistence logic to the repository layer.
type EventsService struct {
	eventsRepository eventsRepository
}

// eventsRepository defines the contract that decouples the service layer
// from the underlying repository logic.
type eventsRepository interface {
	// CreateEvent saves a new event to the repository.
	CreateEvent(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	// DeleteEvent deletes an event from the repository by event ID.
	DeleteEvent(
		ctx context.Context,
		id int,
	) error

	// GetEventByID extracts an Event from the repository by ID, if there is an actual Event with such ID.
	GetEventByID(
		ctx context.Context,
		id int,
	) (domain.Event, error)

	// UpdateEvent patches the given event according to its ID.
	// It uses optimistic concurrency control by checking the event's version
	// to prevent lost updates.
	UpdateEvent(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	// GetEventByTimeRange extracts a slice of Event's according to
	// userID, from and to parameters.
	GetEventByTimeRange(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) []domain.Event
}

// NewEventsService crates a new instance of EventsService.
func NewEventsService(
	eventsRepository eventsRepository,
) *EventsService {
	return &EventsService{
		eventsRepository: eventsRepository,
	}
}
