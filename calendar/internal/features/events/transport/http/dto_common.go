package http

import (
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

// eventDTOresponse is partially included in multiple POST event response operations.
type eventDTOResponse struct {
	ID      int    `json:"id"`
	Version int    `json:"version"`
	Event   string `json:"event"`
	Date    string `json:"date"`
	UserID  int    `json:"user_id"`
}

// resultEventResponse is a response DTO for multiple POST event operations.
type resultEventResponse struct {
	Result eventDTOResponse `json:"result"`
}

func resultEventDTO(domain domain.Event) resultEventResponse {
	eventDTO := eventDTOResponse{
		ID:      domain.ID,
		Version: domain.Version,
		Event:   domain.Event,
		Date:    domain.Date.Format("2006-01-02"),
		UserID:  domain.UserID,
	}

	return resultEventResponse{
		Result: eventDTO,
	}
}

// resultEventsResponse is a response DTO for multiple GET event operations.
type resultEventsResponse struct {
	Result []eventDTOResponse `json:"result"`
}

func resultEventsDTO(domains []domain.Event) resultEventsResponse {
	eventsDTO := make([]eventDTOResponse, 0, len(domains))

	for _, domain := range domains {
		eventDTO := eventDTOResponse{
			ID:      domain.ID,
			Version: domain.Version,
			Event:   domain.Event,
			Date:    domain.Date.Format("2006-01-02"),
			UserID:  domain.UserID,
		}

		eventsDTO = append(eventsDTO, eventDTO)
	}

	return resultEventsResponse{
		Result: eventsDTO,
	}
}
