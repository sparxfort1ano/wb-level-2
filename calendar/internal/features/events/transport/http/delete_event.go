package http

import (
	"fmt"
	"net/http"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/request"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
)

// DeleteEventRequest represents the incoming URL body for deleting an event.
type DeleteEventRequest struct {
	ID int `form:"id" validate:"required,gte=1"`
}

// DeleteEventResponse represents the outgoing JSON body for deleting an event.
type DeleteEventResponse struct {
	Result string `json:"result"`
}

// DeleteEvent processes the HTTP request to delete an event with the given ID.
func (h *EventsHTTPHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctxutil.Logger(ctx)
	responseHandler := response.NewHTTPResponseHandler(w, log)

	var req DeleteEventRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode request and validate HTTP request",
		)
		return
	}

	id := req.ID

	if err := h.eventsService.DeleteEvent(ctx, id); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to delete event",
		)
		return
	}

	response := deleteResultDTO(id)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func deleteResultDTO(id int) DeleteEventResponse {
	return DeleteEventResponse{
		Result: fmt.Sprintf("success with id=%d", id),
	}
}
