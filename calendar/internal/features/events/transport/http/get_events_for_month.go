package http

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
)

// GetEventsForMonthResponse represents the outgoing URL body for getting events for a month.
type GetEventsForMonthResponse resultEventsResponse

// GetEventsForMonth processes an HTTP request to get a list of events for a month according to the userID and date parameters.
func (h *EventsHTTPHandler) GetEventsForMonth(w http.ResponseWriter, r *http.Request) {
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

	eventDomains, err := h.eventsService.GetEventsForMonth(ctx, userID, date)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get events for a month",
		)
		return
	}

	response := GetEventsForMonthResponse(resultEventsDTO(eventDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}
