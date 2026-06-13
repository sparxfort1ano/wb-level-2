package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

func (s *EventsService) CreateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	if err := event.Validate(); err != nil {
		return domain.Event{}, fmt.Errorf("validate event domain: %w", err)
	}

	event, err := s.eventsRepository.CreateEvent(ctx, event)
	if err != nil {
		return domain.Event{}, fmt.Errorf("store new event: %w", err)
	}

	return event, nil
}
