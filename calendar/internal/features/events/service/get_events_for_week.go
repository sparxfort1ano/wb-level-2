package service

import (
	"context"
	"fmt"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (s *EventsService) GetEventsForWeek(
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

	from, to := getFromAndToForWeek(date)

	events := s.eventsRepository.GetEventByTimeRange(ctx, userID, &from, &to)

	return events, nil
}

func getFromAndToForWeek(date *time.Time) (time.Time, time.Time) {
	tmp := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)

	weekDay := int(tmp.Weekday())
	from := tmp.AddDate(0, 0, -((weekDay + 6) % 7))
	to := from.AddDate(0, 0, 7)

	return from, to
}
