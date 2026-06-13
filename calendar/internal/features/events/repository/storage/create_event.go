package storage

import (
	"context"
	"fmt"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (r *EventsRepository) CreateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	eventModel := modelFromDomain(event)

	if err := r.save(&eventModel); err != nil {
		return domain.Event{}, err
	}

	eventDomain := domainFromModel(eventModel)

	return eventDomain, nil
}

// save adds the event to the slice assuming that the user_id exists or has just been created.
// Increments the newUserID field if a request is received without a user_id (new user register).
// Updates the ID and Version of the event model.
func (r *EventsRepository) save(model *EventModel) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if err := r.findOrCreateUserID(model); err != nil {
		return err
	}

	r.updateIDandVersion(model)

	r.storeEvent(*model)

	return nil
}

func (r *EventsRepository) findOrCreateUserID(model *EventModel) error {
	if model.UserID < startUserID || model.UserID > r.newUserID {
		return fmt.Errorf(
			"user with id=%d does not exist: %w",
			model.UserID,
			errs.ErrBadArgument,
		)
	}

	if model.UserID == startUserID {
		r.newUserID++ // user_id's actually >= 1 + startUserID
		model.UserID = r.newUserID
	}

	return nil
}

func (r *EventsRepository) updateIDandVersion(model *EventModel) {
	model.Version = startVersion

	model.ID = r.newEventID
	r.newEventID++
}

func (r *EventsRepository) storeEvent(model EventModel) {
	r.storage[model.ID] = model
}
