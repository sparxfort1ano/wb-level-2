package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (s *EventsService) GetEventsForMonth(
	ctx context.Context,
	userID *int,
	date *time.Time,
) ([]domain.Event, error) {
	if date == nil {
		return nil, fmt.Errorf(
			"required parameter date is omitted: %w",
			errs.ErrInvalidArgument,
		)
	}

	from, to := getFromAndToForMonth(date)

	events := s.eventsRepository.GetEventByTimeRange(ctx, userID, &from, &to)

	return events, nil
}

func getFromAndToForMonth(date *time.Time) (time.Time, time.Time) {
	from := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, time.Local)
	to := time.Date(from.Year(), from.Month()+1, 1, 0, 0, 0, 0, time.Local)

	return from, to
}
