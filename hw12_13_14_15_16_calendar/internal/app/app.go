package app

import (
	"context"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage/memory"
)

type App struct { // TODO
}

type Logger interface { // TODO
}

type Storage interface { // TODO
	Add(ctx context.Context, event memorystorage.EventCreateParams) (string, error) // добавить событие;
	Update(id string, event memorystorage.EventUpdateParams) error                  // обновить событие;
	Delete(id string) error                                                         // удавить событие;
	// получить список событий на день/неделю/месяц;
	GetEventListing(userID string, date time.Time, period string) ([]storage.Event, error)
	Notify(day uint) (string, error)
}

func New(_ Logger, _ Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, _ string, _ string) error { // id title
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
