package app

import (
	"context"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storager interface { // TODO
	AddEventByID(ctx context.Context, event storage.EventCreateDTO, userID string) (string, error)
	UpdateEventByID(ctx context.Context, id string, event storage.EventUpdateDTO, usserID string) error
	DeleteEventByID(ctx context.Context, id string) error
	// получить список событий на день/неделю/месяц;
	GetEventListingByUserID(userID string, date time.Time, period string) ([]storage.Event, error)
	GetEventByID(id string, userID string) (storage.Event, error)
	Notify(day uint) (string, error)
}

func New(_ Logger, _ Storager) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, _ string, _ string) error { // id title
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
