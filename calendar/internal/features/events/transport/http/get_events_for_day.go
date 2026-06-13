package http

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
)

// GetEventsForDayResponse represents the outgoing URL body for getting events for a day.
type GetEventsForDayResponse resultEventsResponse

// GetEventsForDay processes an HTTP request to get a list of events for a day according to the userID and date parameters.
func (h *EventsHTTPHandler) GetEventsForDay(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctxutil.Logger(ctx)
	responseHandler := response.NewHTTPResponseHandler(w, log)

	userID, date, err := getUserIDandDate(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			fmt.Sprintf(
				"failed to get '%s'/'%s'",
				userIDKey,
				dateKey,
			),
		)
		return
	}

	eventDomains, err := h.eventsService.GetEventsForDay(ctx, userID, date)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get events for a day",
		)
		return
	}

	response := GetEventsForDayResponse(resultEventsDTO(eventDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}
