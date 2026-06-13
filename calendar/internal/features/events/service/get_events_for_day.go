package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (s *EventsService) GetEventsForDay(
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

	from, to := getFromAndToForDay(date)

	events := s.eventsRepository.GetEventByTimeRange(ctx, userID, &from, &to)

	return events, nil
}

func getFromAndToForDay(date *time.Time) (time.Time, time.Time) {
	from := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	to := from.AddDate(0, 0, 1)

	return from, to
}
