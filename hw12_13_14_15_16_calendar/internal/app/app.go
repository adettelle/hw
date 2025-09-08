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
	Add(ctx context.Context, event memorystorage.EventToCreate) (string, error) // добавить событие;
	Update(id string, event memorystorage.EventToUpdate) error                  // обновить событие;
	Delete(id string) error                                                     // удавить событие;
	// получить список событий на день/неделю/месяц;
	GetEventListing(userID string, date time.Time, period string) ([]storage.Event, error)
	Notify(day uint) (string, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
