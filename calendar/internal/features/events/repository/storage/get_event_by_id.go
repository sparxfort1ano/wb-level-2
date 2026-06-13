package storage

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (r *EventsRepository) GetEventByID(
	ctx context.Context,
	id int,
) (domain.Event, error) {
	eventModel, err := r.getEventByID(id)
	if err != nil {
		return domain.Event{}, err
	}

	eventDomain := domainFromModel(eventModel)

	return eventDomain, nil
}

func (r *EventsRepository) getEventByID(id int) (EventModel, error) {
	r.rw.RLock()
	defer r.rw.RUnlock()

	if _, ok := r.storage[id]; !ok {
		return EventModel{}, fmt.Errorf(
			"event with id=%d does not exist: %w",
			id,
			errs.ErrBadArgument,
		)
	}

	model := r.storage[id]
	return model, nil
}
