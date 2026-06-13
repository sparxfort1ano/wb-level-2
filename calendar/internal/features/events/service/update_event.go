package service

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

func (s *EventsService) UpdateEvent(
	ctx context.Context,
	id int,
	patch domain.EventPatch,
) (domain.Event, error) {
	event, err := s.eventsRepository.GetEventByID(ctx, id)
	if err != nil {
		return domain.Event{}, fmt.Errorf("get an event: %w", err)
	}

	if err := event.ApplyPatch(patch); err != nil {
		return domain.Event{}, fmt.Errorf("apply event patch: %w", err)
	}

	patchedEvent, err := s.eventsRepository.UpdateEvent(ctx, event)
	if err != nil {
		return domain.Event{}, fmt.Errorf("patch event: %w", err)
	}

	return patchedEvent, nil
}
