package http

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
)

// GetEventsForWeekResponse represents the outgoing URL body for getting events for a week.
type GetEventsForWeekResponse resultEventsResponse

// GetEventsForWeek processes an HTTP request to get a list of events for a week according to the userID and date parameters.
func (h *EventsHTTPHandler) GetEventsForWeek(w http.ResponseWriter, r *http.Request) {
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

	eventDomains, err := h.eventsService.GetEventsForWeek(ctx, userID, date)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get events for a week",
		)
		return
	}

	response := GetEventsForWeekResponse(resultEventsDTO(eventDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}
