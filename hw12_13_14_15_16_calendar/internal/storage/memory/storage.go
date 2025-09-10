package memorystorage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

type Storage struct {
	Events map[string]storage.Event
	mu     sync.RWMutex //nolint:unused
}

func New() *Storage {
	events := map[string]storage.Event{}
	return &Storage{Events: events}
}

type EventCreateParams struct {
	DateStart    time.Time // Дата и время события;
	Title        string
	DateEnd      time.Time // дата и время окончания (Длительность события);
	Description  string    // Описание события - длинный текст, опционально;
	UserID       string    // ID пользователя, владельца события;
	Notification time.Time // За сколько времени высылать уведомление, опционально.
}

type EventUpdateParams struct {
	Date         *time.Time // Дата и время события;
	Title        *string
	Duration     *time.Time // Длительность события (или дата и время окончания); TODO HELP
	Description  *string    // Описание события - длинный текст, опционально;
	Notification *time.Time // За сколько времени высылать уведомление, опционально.  TODO HELP
}

func (s *Storage) Add(_ context.Context, ec EventCreateParams) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := uuid.New().String()
	event := storage.Event{
		ID:           id,
		Title:        ec.Title,
		Date:         ec.DateStart,
		Duration:     ec.DateEnd,
		Description:  ec.Description,
		UserID:       ec.UserID,
		Notification: ec.Notification,
	}
	s.Events[id] = event
	return id, nil
}

func (s *Storage) Update(id string, event EventUpdateParams) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.Events[id]
	if !ok {
		return fmt.Errorf("there is no event with id: %s", id)
	}

	e := s.Events[id]
	if event.Title != nil {
		e.Title = *event.Title
	}
	if event.Date != nil {
		e.Date = *event.Date
	}
	if event.Duration != nil {
		e.Duration = *event.Duration
	}
	if event.Description != nil {
		e.Description = *event.Description
	}
	if event.Notification != nil {
		e.Notification = *event.Notification
	}

	s.Events[id] = e

	return nil
}

func (s *Storage) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.Events[id]
	if !ok {
		return fmt.Errorf("nothing to delete with id: %s", id)
	}
	delete(s.Events, id)
	return nil
}

const (
	day   = "day"
	week  = "week"
	month = "month"
)

// получить список событий на день/неделю/месяц;
// date - дата; если неделя или месяц - то будет вокруг этой даты.
func (s *Storage) GetEventListing(userID string, date time.Time, period string) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var (
		start time.Time
		end   time.Time
	)

	result := []storage.Event{}

	switch period {
	case day:
		start = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 1, 0, time.Local)
		end = time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)
	case week:
		start = StartOfWeek(date)
		end = EndOfWeek(date)
	case month:
		start = StartOfMonth(date)
		end = EndOfMonth(date)
		// default:
		// 	fmt.Errorf("no such period %s", period) // TODO HELP надо ли здесь????
	}

	for _, event := range s.Events {
		if event.Date.After(start) && event.Date.Before(end) && event.UserID == userID {
			result = append(result, event)
		}
	}

	return result, nil
}

func StartOfWeek(date time.Time) time.Time {
	daysSinceSunday := int(date.Weekday())
	s := date.AddDate(0, 0, -daysSinceSunday+1)
	startDate := time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 1, 0, time.Local)
	return startDate
}

func EndOfWeek(date time.Time) time.Time {
	daysUntilSaturday := 7 - int(date.Weekday())
	e := date.AddDate(0, 0, daysUntilSaturday)
	endDate := time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
	return endDate
}

func StartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func EndOfMonth(date time.Time) time.Time {
	firstDayOfNextMonth := StartOfMonth(date).AddDate(0, 1, 0)
	return firstDayOfNextMonth.Add(-time.Second)
}

// получить уведомление за N дней до события.
func (s *Storage) Notify(_ uint) (string, error) { // day
	return "", nil
}

func (s *Storage) GetByID(id string) (storage.Event, error) {
	event, ok := s.Events[id]
	if !ok {
		return storage.Event{}, fmt.Errorf("no evend with id %s", id)
	}
	return event, nil
}
