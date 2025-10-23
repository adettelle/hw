package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/c2fo/testify/require"
)

func TestStorageAdd(t *testing.T) {
	dateStart := time.Now()
	user1 := "1"
	user2 := "2"
	store := New()

	eventToCreate1 := storage.EventCreateDTO{
		Title:       "event1",
		Start:       dateStart,
		End:         dateStart.AddDate(0, 0, 3),
		Description: "description of event 1",
		// UserID:       "1",
		Notification: dateStart.AddDate(0, -1, 0),
	}
	eventToCreate2 := storage.EventCreateDTO{
		Start:       dateStart,
		Title:       "event2",
		End:         dateStart.AddDate(0, 0, 3),
		Description: "description of event 2",
		// UserID:       "2",
		Notification: dateStart.AddDate(0, -2, 0),
	}
	ctx := context.Background()

	id1, err := store.AddEventByID(ctx, eventToCreate1, user1)
	require.NoError(t, err)
	id2, err := store.AddEventByID(ctx, eventToCreate2, user2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)
	event1, err := store.GetEventByID(id1, user1)
	require.NoError(t, err)

	require.Equal(t, eventToCreate1.Start, event1.Start)
	require.Equal(t, eventToCreate1.Title, event1.Title)
	require.Equal(t, eventToCreate1.End, event1.End)
	require.Equal(t, eventToCreate1.Description, event1.Description)
	require.Equal(t, user1, event1.UserID)
	require.Equal(t, eventToCreate1.Notification, event1.Notification)

	event2, err := store.GetEventByID(id2, user2)
	require.NoError(t, err)

	require.Equal(t, eventToCreate2.Start, event2.Start)
	require.Equal(t, eventToCreate2.Title, event2.Title)
	require.Equal(t, eventToCreate2.End, event2.End)
	require.Equal(t, eventToCreate2.Description, event2.Description)
	require.Equal(t, user2, event2.UserID)
	require.Equal(t, eventToCreate2.Notification, event2.Notification)
}

func TestStorageUpdate(t *testing.T) {
	dateStart := time.Now()
	user1 := "1"
	user2 := "2"

	store := New()
	eventToCreate1 := storage.EventCreateDTO{
		Start:       dateStart,
		Title:       "event1",
		End:         dateStart.AddDate(0, 0, 1),
		Description: "description of event 1",
		// UserID:       "1",
		Notification: dateStart.AddDate(0, -1, 0),
	}
	eventToCreate2 := storage.EventCreateDTO{
		Start:       dateStart,
		Title:       "event2",
		End:         dateStart.AddDate(0, 0, 3),
		Description: "description of event 2",
		// UserID:       "2",
		Notification: dateStart.AddDate(0, -2, 0),
	}

	ctx := context.Background()
	_, err := store.AddEventByID(ctx, eventToCreate1, user1)
	require.NoError(t, err)
	id2, err := store.AddEventByID(ctx, eventToCreate2, user2)
	require.NoError(t, err)

	title := "new event2"
	description := "new description of event 2"
	eventToUpdate2 := storage.EventUpdateDTO{
		Title:       &title,
		Description: &description,
	}

	err = store.UpdateEventByID(ctx, id2, eventToUpdate2, user2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)

	event2, err := store.GetEventByID(id2, user2)
	require.NoError(t, err)

	require.Equal(t, eventToCreate2.Start, event2.Start)
	require.Equal(t, *eventToUpdate2.Title, event2.Title)
	require.Equal(t, eventToCreate2.End, event2.End)
	require.Equal(t, *eventToUpdate2.Description, event2.Description)
	require.Equal(t, user2, event2.UserID)
	require.Equal(t, eventToCreate2.Notification, event2.Notification)
}

func TestStorageDelete(t *testing.T) {
	dateStart := time.Now()
	user1 := "1"
	user2 := "2"

	store := New()
	eventToCreate1 := storage.EventCreateDTO{
		Start:       dateStart,
		Title:       "event1",
		End:         dateStart.AddDate(0, 0, 3),
		Description: "description of event 1",
		// UserID:       "1",
		Notification: dateStart.AddDate(0, -1, 0),
	}
	eventToCreate2 := storage.EventCreateDTO{
		Start:       dateStart,
		Title:       "event2",
		End:         dateStart.AddDate(0, 0, 3),
		Description: "description of event 2",
		// UserID:       "2",
		Notification: dateStart.AddDate(0, -2, 0),
	}
	ctx := context.Background()
	id1, err := store.AddEventByID(ctx, eventToCreate1, user1)
	require.NoError(t, err)
	id2, err := store.AddEventByID(ctx, eventToCreate2, user2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)
	_, err = store.GetEventByID(id1, user1)
	require.NoError(t, err)

	store.DeleteEventByID(ctx, id1)
	require.Equal(t, len(store.Events), 1)

	event2, err := store.GetEventByID(id2, user2)
	require.NoError(t, err)

	require.Equal(t, eventToCreate2.Start, event2.Start)
	require.Equal(t, eventToCreate2.Title, event2.Title)
	require.Equal(t, eventToCreate2.End, event2.End)
	require.Equal(t, eventToCreate2.Description, event2.Description)
	require.Equal(t, user2, event2.UserID)
	require.Equal(t, eventToCreate2.Notification, event2.Notification)
}

func TestStorageGetEventListing(t *testing.T) {
	store := New()
	date1 := time.Date(2025, time.September, 4, 10, 0, 0, 0, time.Local)
	date2 := time.Date(2025, time.September, 5, 10, 0, 0, 0, time.Local)
	date3 := time.Date(2025, time.September, 20, 10, 0, 0, 0, time.Local)
	user1 := "1"

	eventToCreate1 := storage.EventCreateDTO{
		Start:       date1,
		Title:       "event1",
		End:         date1.AddDate(0, 0, 1),
		Description: "description of event 1",
		// UserID:       "1",
		Notification: date1.AddDate(0, -1, 0),
	}
	eventToCreate2 := storage.EventCreateDTO{
		Start:       date2,
		Title:       "event2",
		End:         date2.AddDate(0, 0, 2),
		Description: "description of event 2",
		// UserID:       "1",
		Notification: date2.AddDate(0, -2, 0),
	}
	eventToCreate3 := storage.EventCreateDTO{
		Start:       date3,
		Title:       "event3",
		End:         date3.AddDate(0, 0, 3),
		Description: "description of event 3",
		// UserID:       "1",
		Notification: date3.AddDate(0, -2, 0),
	}
	ctx := context.Background() // TODO HELP

	_, err := store.AddEventByID(ctx, eventToCreate1, user1)
	require.NoError(t, err)
	id2, err := store.AddEventByID(ctx, eventToCreate2, user1)
	require.NoError(t, err)
	_, err = store.AddEventByID(ctx, eventToCreate3, user1)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 3)

	res1, err := store.GetEventListingByUserID("1", date2, day)
	require.NoError(t, err)
	require.Equal(t, len(res1), 1)
	event2, err := store.GetEventByID(id2, user1)
	require.NoError(t, err)
	require.Equal(t, res1[0], event2)

	res2, err := store.GetEventListingByUserID("1", date1, week)
	require.NoError(t, err)
	require.Equal(t, len(res2), 2)

	for _, event := range res2 {
		_, ok := store.Events[event.ID]
		require.True(t, ok)
	}

	res3, err := store.GetEventListingByUserID("1", date1, month)
	require.NoError(t, err)
	require.Equal(t, len(res3), 3)

	for _, event := range res3 {
		_, ok := store.Events[event.ID]
		require.True(t, ok)
	}
}
