package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/c2fo/testify/require"
)

func TestStorageAdd(t *testing.T) {
	dateStart := time.Now()

	store := New()
	eventToCreate1 := EventCreateParams{
		DateStart:    dateStart,
		Title:        "event1",
		DateEnd:      dateStart.AddDate(0, 0, 3),
		Description:  "description of event 1",
		UserID:       "1",
		Notification: dateStart.AddDate(0, -1, 0),
	}
	eventToCreate2 := EventCreateParams{
		DateStart:    dateStart,
		Title:        "event2",
		DateEnd:      dateStart.AddDate(0, 0, 3),
		Description:  "description of event 2",
		UserID:       "2",
		Notification: dateStart.AddDate(0, -2, 0),
	}
	ctx := context.Background() // TODO HELP вот так просто сделать здесь контекст или как-то иначе ?????

	id1, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)
	event1, err := store.GetByID(id1)
	require.NoError(t, err)
	// require.Equal(t, eventToCreate1, event1)

	require.Equal(t, eventToCreate1.DateStart, event1.Date)
	require.Equal(t, eventToCreate1.Title, event1.Title)
	require.Equal(t, eventToCreate1.DateEnd, event1.Duration)
	require.Equal(t, eventToCreate1.Description, event1.Description)
	require.Equal(t, eventToCreate1.UserID, event1.UserID)
	require.Equal(t, eventToCreate1.Notification, event1.Notification)

	event2, err := store.GetByID(id2)
	require.NoError(t, err)
	// require.Equal(t, eventToCreate2, event2)
	require.Equal(t, eventToCreate2.DateStart, event2.Date)
	require.Equal(t, eventToCreate2.Title, event2.Title)
	require.Equal(t, eventToCreate2.DateEnd, event2.Duration)
	require.Equal(t, eventToCreate2.Description, event2.Description)
	require.Equal(t, eventToCreate2.UserID, event2.UserID)
	require.Equal(t, eventToCreate2.Notification, event2.Notification)
}

func TestStorageUpdate(t *testing.T) {
	dateStart := time.Now()

	store := New()
	eventToCreate1 := EventCreateParams{
		DateStart:    dateStart,
		Title:        "event1",
		DateEnd:      dateStart.AddDate(0, 0, 1),
		Description:  "description of event 1",
		UserID:       "1",
		Notification: dateStart.AddDate(0, -1, 0),
	}
	eventToCreate2 := EventCreateParams{
		DateStart:    dateStart,
		Title:        "event2",
		DateEnd:      dateStart.AddDate(0, 0, 3),
		Description:  "description of event 2",
		UserID:       "2",
		Notification: dateStart.AddDate(0, -2, 0),
	}

	ctx := context.Background() // TODO HELP
	_, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)

	title := "new event2"
	description := "new description of event 2"
	eventToUpdate2 := EventUpdateParams{
		Title:       &title,
		Description: &description,
	}

	err = store.Update(id2, eventToUpdate2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)

	event2, err := store.GetByID(id2)
	require.NoError(t, err)

	require.Equal(t, eventToCreate2.DateStart, event2.Date)
	require.Equal(t, *eventToUpdate2.Title, event2.Title)
	require.Equal(t, eventToCreate2.DateEnd, event2.Duration)
	require.Equal(t, *eventToUpdate2.Description, event2.Description)
	require.Equal(t, eventToCreate2.UserID, event2.UserID)
	require.Equal(t, eventToCreate2.Notification, event2.Notification)
}

func TestStorageDelete(t *testing.T) {
	dateStart := time.Now()

	store := New()
	eventToCreate1 := EventCreateParams{
		DateStart:    dateStart,
		Title:        "event1",
		DateEnd:      dateStart.AddDate(0, 0, 3),
		Description:  "description of event 1",
		UserID:       "1",
		Notification: dateStart.AddDate(0, -1, 0),
	}
	eventToCreate2 := EventCreateParams{
		DateStart:    dateStart,
		Title:        "event2",
		DateEnd:      dateStart.AddDate(0, 0, 3),
		Description:  "description of event 2",
		UserID:       "2",
		Notification: dateStart.AddDate(0, -2, 0),
	}
	ctx := context.Background() // TODO HELP
	id1, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)
	_, err = store.GetByID(id1)
	require.NoError(t, err)

	store.Delete(id1)
	require.Equal(t, len(store.Events), 1)

	event2, err := store.GetByID(id2)
	require.NoError(t, err)

	require.Equal(t, eventToCreate2.DateStart, event2.Date)
	require.Equal(t, eventToCreate2.Title, event2.Title)
	require.Equal(t, eventToCreate2.DateEnd, event2.Duration)
	require.Equal(t, eventToCreate2.Description, event2.Description)
	require.Equal(t, eventToCreate2.UserID, event2.UserID)
	require.Equal(t, eventToCreate2.Notification, event2.Notification)
}

func TestStorageGetEventListing(t *testing.T) {
	store := New()
	date1 := time.Date(2025, time.September, 4, 10, 0, 0, 0, time.Local)
	date2 := time.Date(2025, time.September, 5, 10, 0, 0, 0, time.Local)
	date3 := time.Date(2025, time.September, 20, 10, 0, 0, 0, time.Local)

	eventToCreate1 := EventCreateParams{
		DateStart:    date1,
		Title:        "event1",
		DateEnd:      date1.AddDate(0, 0, 1),
		Description:  "description of event 1",
		UserID:       "1",
		Notification: date1.AddDate(0, -1, 0),
	}
	eventToCreate2 := EventCreateParams{
		DateStart:    date2,
		Title:        "event2",
		DateEnd:      date2.AddDate(0, 0, 2),
		Description:  "description of event 2",
		UserID:       "1",
		Notification: date2.AddDate(0, -2, 0),
	}
	eventToCreate3 := EventCreateParams{
		DateStart:    date3,
		Title:        "event3",
		DateEnd:      date3.AddDate(0, 0, 3),
		Description:  "description of event 3",
		UserID:       "1",
		Notification: date3.AddDate(0, -2, 0),
	}
	ctx := context.Background() // TODO HELP

	_, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)
	_, err = store.Add(ctx, eventToCreate3)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 3)

	res1, err := store.GetEventListing("1", date2, day)
	require.NoError(t, err)
	require.Equal(t, len(res1), 1)
	event2, err := store.GetByID(id2)
	require.NoError(t, err)
	require.Equal(t, res1[0], event2)

	res2, err := store.GetEventListing("1", date1, week)
	require.NoError(t, err)
	require.Equal(t, len(res2), 2)

	for _, event := range res2 {
		_, ok := store.Events[event.ID]
		require.True(t, ok)
	}

	res3, err := store.GetEventListing("1", date1, month)
	require.NoError(t, err)
	require.Equal(t, len(res3), 3)

	for _, event := range res3 {
		_, ok := store.Events[event.ID]
		require.True(t, ok)
	}
}
