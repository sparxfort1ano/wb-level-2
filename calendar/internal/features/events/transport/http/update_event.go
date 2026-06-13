package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/ctxutil"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
	errs "github.com/sparxfort1ano/wb-level-2/calendar/internal/core/errors"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/request"
	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/response"
)

// UpdateEventRequest represents the incoming URL body for updating an event.
type UpdateEventRequest struct {
	ID    int     `form:"id" validate:"required,gte=1"`
	Event *string `form:"event" validate:"omitempty,min=1"`
	Date  *string `form:"date" validate:"omitempty"`
}

// UpdateEventResponse represents the outgoing JSON body for updating an event.
type UpdateEventResponse resultEventResponse

// UpdateEvent processes the HTTP request to partially update an existing event with the given ID.
func (h *EventsHTTPHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := ctxutil.Logger(ctx)
	responseHandler := response.NewHTTPResponseHandler(w, log)

	var req UpdateEventRequest
	if err := request.DecodeAndValidateRequest(r, &req); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)
		return
	}

	eventPatchDomain, err := domainPatchFromDTO(req)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'data'",
		)
		return
	}

	id := req.ID

	eventDomain, err := h.eventsService.UpdateEvent(ctx, id, eventPatchDomain)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to update an event",
		)
		return
	}

	response := UpdateEventResponse(resultEventDTO(eventDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func domainPatchFromDTO(dto UpdateEventRequest) (domain.EventPatch, error) {
	var date *time.Time

	if dto.Date != nil {
		layout := "2006-01-02"
		tmp, err := time.Parse(layout, *dto.Date)
		if err != nil {
			return domain.EventPatch{}, fmt.Errorf(
				"value %s is not a valid date: %v: %w",
				*dto.Date,
				err,
				errs.ErrInvalidArgument,
			)
		}
		date = &tmp
	}

	return *domain.NewEventPatch(
		dto.Event,
		date,
	), nil
}
