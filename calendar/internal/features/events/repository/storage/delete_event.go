package storage

import (
	"context"
	"fmt"

	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
)

func (r *EventsRepository) DeleteEvent(
	ctx context.Context,
	id int,
) error {
	if err := r.deleteEvent(id); err != nil {
		return err
	}

	return nil
}

func (r *EventsRepository) deleteEvent(id int) error {
	r.rw.Lock()
	defer r.rw.Unlock()

	if _, ok := r.storage[id]; !ok {
		return fmt.Errorf(
			"event with id=%d does not exist: %w",
			id,
			errs.ErrBadArgument,
		)
	}

	delete(r.storage, id)

	return nil
}
