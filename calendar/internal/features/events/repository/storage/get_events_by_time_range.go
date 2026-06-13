package storage

import (
	"context"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

func (r *EventsRepository) GetEventByTimeRange(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) []domain.Event {
	eventModels := r.getEventByTimeRange(userID, from, to)

	eventDomains := domainsFromModels(eventModels)

	return eventDomains
}

func (r *EventsRepository) getEventByTimeRange(userID *int, from *time.Time, to *time.Time) []EventModel {
	models := make([]EventModel, 0)

	r.rw.RLock()
	defer r.rw.RUnlock()

	for _, model := range r.storage {
		if userID == nil || *userID == model.UserID {
			date := model.Date
			if (from == nil || !date.Before(*from)) && (to == nil || date.Before(*to)) {
				models = append(models, model)
			}
		}
	}

	return models
}
