package storage

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (r *EventsRepository) UpdateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	eventModel := modelFromDomain(event)

	patchedEventModel, err := r.updateEvent(eventModel)
	if err != nil {
		return domain.Event{}, err
	}

	patchedEventDomain := domainFromModel(patchedEventModel)

	return patchedEventDomain, nil
}

func (r *EventsRepository) updateEvent(patchedModel EventModel) (EventModel, error) {
	id := patchedModel.ID

	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.storage[id]; !ok {
		return EventModel{}, fmt.Errorf(
			"event with id=%d deleted while patching: %w",
			id,
			errs.ErrBadArgument,
		)
	}

	prevVersion := r.storage[id].Version
	if prevVersion != patchedModel.Version {
		return EventModel{}, fmt.Errorf(
			"event id=%d version mismatch (before processing update:%d, after processing update:%d): %w",
			id,
			patchedModel.Version,
			prevVersion,
			errs.ErrBadArgument,
		)
	}

	patchedModel.Version++

	r.storage[id] = patchedModel

	return patchedModel, nil
}
