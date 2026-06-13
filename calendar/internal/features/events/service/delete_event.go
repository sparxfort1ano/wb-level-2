package service

import (
	"context"
	"fmt"
)

func (s *EventsService) DeleteEvent(
	ctx context.Context,
	id int,
) error {
	if err := s.eventsRepository.DeleteEvent(ctx, id); err != nil {
		return fmt.Errorf("delete new event from storage: %w", err)
	}

	return nil
}
