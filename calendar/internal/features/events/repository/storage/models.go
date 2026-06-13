package storage

import (
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

// EventModel represents the storage schema for an event.
type EventModel struct {
	ID      int
	Version int
	Event   string
	Date    time.Time
	UserID  int
}

func modelFromDomain(domain domain.Event) EventModel {
	return EventModel{
		ID:      domain.ID,
		Version: domain.Version,
		Event:   domain.Event,
		Date:    domain.Date,
		UserID:  domain.UserID,
	}
}

func domainFromModel(model EventModel) domain.Event {
	return domain.NewEvent(
		model.ID,
		model.Version,
		model.Event,
		model.Date,
		model.UserID,
	)
}

func domainsFromModels(models []EventModel) []domain.Event {
	domains := make([]domain.Event, 0, len(models))

	for _, model := range models {
		domain := domainFromModel(model)

		domains = append(domains, domain)
	}

	return domains
}
