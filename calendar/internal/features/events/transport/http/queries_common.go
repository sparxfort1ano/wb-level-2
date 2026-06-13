package http

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/transport/http/request"
)

const (
	userIDKey = "user_id"
	dateKey   = "date"
)

func getUserIDandDate(r *http.Request) (*int, *time.Time, error) {
	userID, errUserID := request.GetIntQueryParam(r, userIDKey)
	date, errDate := request.GetDateQueryParam(r, dateKey)

	if errs := errors.Join(
		errUserID,
		errDate,
	); errs != nil {
		return nil, nil, fmt.Errorf("get query params: %w", errs)
	}

	return userID, date, nil
}
