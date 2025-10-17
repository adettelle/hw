package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
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
	Start        time.Time // Дата и время события;
	End          time.Time // дата и время окончания (Длительность события);
	Description  string    // Описание события - длинный текст, опционально;
	Notification time.Time
	// (дата и время, когда высылать уведомление) За сколько времени высылать уведомление, опционально
	Notified bool
}

func (s *DBStorage) GetEventByID(eventID string, userID string) (storage.Event, error) {
	sqlSt := `SELECT title, created_at, date_start, date_end, description, notification, notified
	 	FROM event WHERE account_id = $1 and id = $2;`
	row := s.DB.QueryRowContext(s.Ctx, sqlSt, userID, eventID)

	var e eventGetByID

	err := row.Scan(&e.Title, &e.CreatedAt, &e.Start, &e.End,
		&e.Description, &e.Notification, &e.Notified)
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
		Start:        e.Start,
		End:          e.End,
		Description:  e.Description,
		UserID:       userID,
		Notification: e.Notification,
		Notified:     e.Notified,
	}
	return event, err
}

func (s *DBStorage) AddEventByID(ctx context.Context,
	e storage.EventCreateDTO, userID string,
) (string, error) { // user_id,
	sqlSt := `insert into event (title, date_start, date_end, 
		description, account_id, notification, notified) 
		values ($1, $2, $3, $4, $5, $6, $7) returning id;`

	row := s.DB.QueryRowContext(ctx, sqlSt, e.Title, e.Start,
		e.End, e.Description, userID, e.Notification, e.Notified)

	var eventID string
	err := row.Scan(&eventID)
	if err != nil {
		return "", err
	}

	s.Logg.Info("Event have been added")

	return eventID, nil
}

func (s *DBStorage) UpdateEventByID(ctx context.Context,
	eventID string, event storage.EventUpdateDTO, userID string,
) error {
	pairs := map[string]any{}

	if event.Title != nil {
		pairs["title"] = event.Title
	}
	if event.Start != nil {
		pairs["date_start"] = event.Start
	}
	if event.End != nil {
		pairs["date_end"] = event.End
	}
	if event.Description != nil {
		pairs["description"] = event.Description
	}
	if event.Notification != nil {
		pairs["notification"] = event.Notification
	}
	if !event.Notified {
		pairs["notified"] = event.Notified
	}

	sqlStBase := `update event set `
	sqlSet := []string{}
	vals := []any{} // "a", "b"
	// pairs = title: a, description: b
	index := 1

	for k, v := range pairs {
		sqlSet = append(sqlSet, fmt.Sprintf("%s = $%d", k, index))
		index++
		vals = append(vals, v)
	}
	// sqlSet = ["title = $1", "description = $2"]
	// vals = ["a", "b"]

	if len(pairs) == 0 {
		s.Logg.Info("no field to update", zap.String("eventID", eventID))
		return nil
	}
	vals = append(vals, eventID, userID)

	sqlSt := sqlStBase + strings.Join(sqlSet, ", ") +
		fmt.Sprintf(" where id = $%d and account_id = $%d;", index, index+1)
		// " where id = $" + strconv.Itoa(index) +
		// " and account_id = $" + strconv.Itoa(index+1) + ";" //nolint:gosec
	// update event set title = $1, description = $2 where id = $3 and account_id = $4;
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
	fmt.Println(" !!!!!!!! date:", date)
	var sqlSt string

	switch period {
	case day:
		sqlSt = `SELECT id, title, date_start, date_end, description, notification, notified
			FROM event
			WHERE account_id = $1
			AND date_start >= $2::date
			AND date_start <  ($2::date + INTERVAL '1 day')
			ORDER BY date_start;`
	case week:
		sqlSt = `SELECT id, title, date_start, date_end, description, notification, notified
			FROM event
			WHERE account_id = $1
			AND date_start >= $2::date
			AND date_start <  ($2::date + INTERVAL '1 week')
			ORDER BY date_start;`
	case month:
		sqlSt = `SELECT id, title, date_start, date_end, description, notification, notified
			FROM event
			WHERE account_id = $1
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
		err := rows.Scan(&e.ID, &e.Title, &e.Start, &e.End,
			&e.Description, &e.Notification, &e.Notified)
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

func (s *DBStorage) SetNotified(ctx context.Context, ids []string) ([]string, error) {
	if len(ids) == 0 {
		s.Logg.Info("nothing to notify.")
		return nil, nil
	}
	s.Logg.Info("setting notified events.", zap.Int("amount", len(ids)))

	sqlBase := `update event set notified = true where account_id = 1 `
	sqlSt := sqlBase + fmt.Sprintf("and id in (%s) returning id;", strings.Join(ids, ", "))
	rows, err := s.DB.QueryContext(ctx, sqlSt)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil || rows.Err() != nil {
			return nil, err
		}
		result = append(result, id)
	}
	s.Logg.Info("set notified events.")
	return result, nil
}

func (s *DBStorage) CollectEventsToNotify(ctx context.Context) ([]storage.EventToNotify, error) {
	s.Logg.Info("collecting events to notify.")

	var events []storage.EventToNotify

	sqlSt := `SELECT id, title, date_start, account_id 
		from event where date_start 
		between now() and (now() + interval '1' hour) and notified = false;`

	rows, err := s.DB.QueryContext(ctx, sqlSt)
	if err != nil || rows.Err() != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e storage.EventToNotify
		err := rows.Scan(&e.ID, &e.Title, &e.Start, &e.UserID)
		if err != nil || rows.Err() != nil {
			return nil, err
		}
		events = append(events, e)
	}
	s.Logg.Info("events to notify are collected.")

	return events, nil
}

func (s *DBStorage) DeleteEvents(ctx context.Context) error {
	s.Logg.Info("cleaning outdated events.")
	sqlSt := `delete from event 
		where account_id = 1 and date_end < (now() - interval '1' year);`

	_, err := s.DB.ExecContext(ctx, sqlSt)
	if err != nil {
		s.Logg.Error("error in deleting event from DB", zap.Error(err))
		return err
	}

	s.Logg.Info("DB is clean.")
	return nil
}
