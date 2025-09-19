package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"go.uber.org/zap"
)

type DBStorage struct {
	Ctx  context.Context
	DB   *sql.DB
	Logg *zap.Logger
}

func New(ctx context.Context, db *sql.DB, logg *zap.Logger) *DBStorage {
	return &DBStorage{Ctx: ctx, DB: db, Logg: logg}
}

type eventGetByID struct {
	Title        string
	CreatedAt    time.Time
	Date         time.Time // Дата и время события;
	Duration     time.Time // дата и время окончания (Длительность события);
	Description  string    // Описание события - длинный текст, опционально;
	Notification time.Time
	// (дата и время, когда высылать уведомление) За сколько времени высылать уведомление, опционально.
}

func (s *DBStorage) GetEventByID(eventID string, userID string) (storage.Event, error) {
	sqlSt := `SELECT title, created_at, date_start, date_end, description, notification
	 	FROM event WHERE account_id = $1 and id = $2;`
	row := s.DB.QueryRowContext(s.Ctx, sqlSt, userID, eventID)

	var e eventGetByID

	err := row.Scan(&e.Title, &e.CreatedAt, &e.Date, &e.Duration, &e.Description, &e.Notification)
	if err != nil {
		if err == sql.ErrNoRows {
			s.Logg.Error("no event in DB", zap.Error(err), zap.String("eventID", eventID))
			return storage.Event{}, err
		}
		s.Logg.Error("error in getting event by id", zap.Error(err), zap.String("eventID", eventID))
		return storage.Event{}, err
	}

	event := storage.Event{
		ID:           eventID,
		Title:        e.Title,
		CreatedAt:    e.CreatedAt,
		Date:         e.Date,
		Duration:     e.Duration,
		Description:  e.Description,
		UserID:       userID,
		Notification: e.Notification,
	}
	return event, err
}

func (s *DBStorage) AddEventByID(ctx context.Context,
	e storage.EventCreateDTO, userID string,
) (string, error) { // user_id,
	sqlSt := `insert into event (title, date_start, date_end, description, account_id, notification) 
		values ($1, $2, $3, $4, $5, $6) returning id;`

	row := s.DB.QueryRowContext(ctx, sqlSt, e.Title, e.DateStart,
		e.DateEnd, e.Description, userID, e.Notification)

	var eventID string
	err := row.Scan(&eventID)
	if err != nil {
		return "", err
	}

	s.Logg.Info("Event have been added")

	return eventID, nil
}

func (s *DBStorage) UpdateEventByID(ctx context.Context,
	eventID string, event storage.EventUpdateDTO, _ string,
) error { // userID
	pairs := map[string]any{}

	if event.Title != nil {
		pairs["title"] = event.Title
	}
	if event.Date != nil {
		pairs["date_start"] = event.Date
	}
	if event.Duration != nil {
		pairs["date_end"] = event.Duration
	}
	if event.Description != nil {
		pairs["description"] = event.Description
	}
	if event.Notification != nil {
		pairs["notification"] = event.Notification
	}

	sqlStBase := `update event set `
	sqlSet := []string{}
	vals := []any{}
	index := 1

	for k, v := range pairs {
		sqlSet = append(sqlSet, fmt.Sprintf("%s = $%d", k, index))
		index++
		vals = append(vals, v)
	}

	if len(pairs) == 0 {
		s.Logg.Info("no field to update", zap.String("eventID", eventID))
		return nil
	}
	vals = append(vals, eventID)

	sqlSt := sqlStBase + strings.Join(sqlSet, ", ") + " where id = $" + strconv.Itoa(index) + ";" //nolint:gosec

	_, err := s.DB.ExecContext(ctx, sqlSt, vals...)
	if err != nil {
		s.Logg.Error("error in updateing event", zap.Error(err), zap.String("eventID", eventID))
		return err
	}

	return nil
}

func (s *DBStorage) DeleteEventByID(ctx context.Context, eventID string) error {
	sqlSt := `delete from event where id = $1;`

	_, err := s.DB.ExecContext(ctx, sqlSt, eventID)
	if err != nil {
		s.Logg.Error("error in deleting event from DB", zap.Error(err), zap.String("eventID", eventID))
		return err
	}

	s.Logg.Info("Event is deleted.")
	return nil
}

const (
	day   = "day"
	week  = "week"
	month = "month"
)

// получить список событий на день/неделю/месяц.
func (s *DBStorage) GetEventListingByUserID(userID string, date time.Time, period string) ([]storage.Event, error) {
	events := []storage.Event{}

	var sqlSt string

	switch period {
	case day:
		sqlSt = `SELECT id, title, date_start, date_end, description, notification
			FROM event
			WHERE account_id = $1
			AND date_start >= $2::date
			AND date_start <  ($2::date + INTERVAL '1 day')
			ORDER BY date_start;`
	case week:
		sqlSt = `SELECT id, title, date_start, date_end, description, notification
			FROM event
			WHERE account_id = 1
			AND date_start >= $2::date
			AND date_start <  ($2::date + INTERVAL '1 week')
			ORDER BY date_start;`
	case month:
		sqlSt = `SELECT id, title, date_start, date_end, description, notification
			FROM event
			WHERE account_id = 1
			AND date_start >= $2::date
			AND date_start <  ($2::date + INTERVAL '1 month')
			ORDER BY date_start;`
	}

	rows, err := s.DB.QueryContext(context.Background(), sqlSt, userID, date)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.Event
		err := rows.Scan(&e.ID, &e.Title, &e.Date, &e.Duration, &e.Description, &e.Notification)
		if err != nil || rows.Err() != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (s *DBStorage) Notify(_ uint) (string, error) { // day
	return "", nil
}

func (s *DBStorage) Connect(_ context.Context) error {
	// TODO
	return nil
}

func (s *DBStorage) Close(_ context.Context) error {
	// TODO
	return nil
}
