package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sparxfort1ano/wb-level-2/calendar/internal/core/domain"
)

type mockEventRepository struct {
	CreateEventFunc func(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	DeleteEventFunc func(
		ctx context.Context,
		id int,
	) error

	GetIventByIDFunc func(
		ctx context.Context,
		id int,
	) (domain.Event, error)

	UpdateEventFunc func(
		ctx context.Context,
		event domain.Event,
	) (domain.Event, error)

	GetEventByTimeRangeFunc func(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) []domain.Event
}

func (m *mockEventRepository) CreateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	return m.CreateEventFunc(ctx, event)
}

func (m *mockEventRepository) DeleteEvent(
	ctx context.Context,
	id int,
) error {
	return m.DeleteEventFunc(ctx, id)
}

func (m *mockEventRepository) GetEventByID(
	ctx context.Context,
	id int,
) (domain.Event, error) {
	return m.GetIventByIDFunc(ctx, id)
}

func (m *mockEventRepository) UpdateEvent(
	ctx context.Context,
	event domain.Event,
) (domain.Event, error) {
	return m.UpdateEventFunc(ctx, event)
}

func (m *mockEventRepository) GetEventByTimeRange(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) []domain.Event {
	return m.GetEventByTimeRangeFunc(ctx, userID, from, to)
}

func TestCreateEvent(t *testing.T) {
	now := time.Now()
	past := now.AddDate(0, 0, -1)

	testCases := []struct {
		name          string
		inputEvent    domain.Event
		mockFunc      func(m *mockEventRepository)
		isErrExpected bool
	}{
		{
			name:       "success:valid domain.Event fields",
			inputEvent: domain.NewEventUnitialized("football", now, 1),
			mockFunc: func(m *mockEventRepository) {
				m.CreateEventFunc = func(ctx context.Context, event domain.Event) (domain.Event, error) {
					event.ID = 2
					event.Version = 1
					return event, nil
				}
			},
			isErrExpected: false,
		},
		{
			name:       "invalid argument:empty event string",
			inputEvent: domain.NewEventUnitialized("", now, 1),
			mockFunc: func(m *mockEventRepository) {
				m.CreateEventFunc = func(ctx context.Context, event domain.Event) (domain.Event, error) {
					t.Fatal("repository should not be called: validation error")
					return domain.Event{}, nil
				}
			},
			isErrExpected: true,
		},
		{
			name:       "invalid argument:date from past",
			inputEvent: domain.NewEventUnitialized("not empty", past, 1),
			mockFunc: func(m *mockEventRepository) {
				m.CreateEventFunc = func(ctx context.Context, event domain.Event) (domain.Event, error) {
					t.Fatal("repository should not be called: validation error")
					return domain.Event{}, nil
				}
			},
			isErrExpected: true,
		},
		{
			name:       "error from repository layer",
			inputEvent: domain.NewEventUnitialized("some big data", now, 1),
			mockFunc: func(m *mockEventRepository) {
				m.CreateEventFunc = func(ctx context.Context, event domain.Event) (domain.Event, error) {
					return domain.Event{}, errors.New("unknown error")
				}
			},
			isErrExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockEventRepository{}

			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo)
			}

			service := NewEventsService(mockRepo)

			_, err := service.CreateEvent(context.Background(), tc.inputEvent)
			if (err != nil) != tc.isErrExpected {
				t.Errorf("create_event: error=%v, error expectation=%t", err, tc.isErrExpected)
			}
		})
	}
}

func TestGetEventsForDay(t *testing.T) {
	randomUserID := 1

	randomDate := time.Date(2026, time.December, 31, 15, 16, 17, 18, time.Local)
	expectedFrom := time.Date(2026, time.December, 31, 0, 0, 0, 0, time.Local)
	expectedTo := time.Date(2027, time.January, 1, 0, 0, 0, 0, time.Local)

	testCases := []struct {
		name           string
		inputUserID    *int
		inputDate      *time.Time
		mockFunc       func(m *mockEventRepository, t *testing.T)
		expectedEvents int
		isErrExpected  bool
	}{
		{
			name:        "success:valid from and to calculating",
			inputUserID: &randomUserID,
			inputDate:   &randomDate,
			mockFunc: func(m *mockEventRepository, t *testing.T) {
				m.GetEventByTimeRangeFunc = func(ctx context.Context, userID *int, from, to *time.Time) []domain.Event {
					if !from.Equal(expectedFrom) {
						t.Errorf("expected from=%v, got from=%v", expectedFrom, from)
					}

					if !to.Equal(expectedTo) {
						t.Errorf("expected to=%v, got to=%v", expectedTo, to)
					}

					return []domain.Event{
						{
							ID:    1,
							Event: "random event 1",
						},
						{
							ID:    2,
							Event: "random event 2",
						},
					}
				}
			},
			expectedEvents: 2,
			isErrExpected:  false,
		},
		{
			name:        "invalid argument:omitted date",
			inputUserID: &randomUserID,
			inputDate:   nil,
			mockFunc: func(m *mockEventRepository, t *testing.T) {
				m.GetEventByTimeRangeFunc = func(ctx context.Context, userID *int, from, to *time.Time) []domain.Event {
					t.Fatal("repository should not be called: validation error")
					return nil
				}
			},
			expectedEvents: 0,
			isErrExpected:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mockEventRepository{}

			if tc.mockFunc != nil {
				tc.mockFunc(mockRepo, t)
			}

			service := NewEventsService(mockRepo)

			events, err := service.GetEventsForDay(context.Background(), tc.inputUserID, tc.inputDate)
			if (err != nil) != tc.isErrExpected {
				t.Errorf("get_events_for_day: %v", err)
			}

			if len(events) != tc.expectedEvents {
				t.Errorf("get_events_for_day: expected events=%d, got events=%d", tc.expectedEvents, len(events))
			}
		})
	}
}
