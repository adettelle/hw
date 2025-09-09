package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/c2fo/testify/require"
)

func TestStorageAdd(t *testing.T) {
	store := New()
	eventToCreate1 := EventToCreate{
		DateStart:    time.Now(),
		Title:        "event1",
		DateEnd:      3,
		Description:  "description of event 1",
		UserID:       "1",
		Notification: time.Now().AddDate(0, -1, 0),
	}
	eventToCreate2 := EventToCreate{
		DateStart:    time.Now(),
		Title:        "event2",
		DateEnd:      3,
		Description:  "description of event 2",
		UserID:       "2",
		Notification: time.Now().AddDate(0, -2, 0),
	}
	ctx := context.Background() // TODO HELP вот так просто сделать здесь контекст или как-то иначе ?????

	id1, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)
	_, ok := store.Events[id1]
	require.True(t, ok)
	_, ok = store.Events[id2]
	require.True(t, ok)
}

func TestStorageUpdate(t *testing.T) {
	store := New()
	eventToCreate1 := EventToCreate{
		DateStart:    time.Now(),
		Title:        "event1",
		DateEnd:      1,
		Description:  "description of event 1",
		UserID:       "1",
		Notification: time.Now().AddDate(0, -1, 0),
	}
	eventToCreate2 := EventToCreate{
		DateStart:    time.Now(),
		Title:        "event2",
		DateEnd:      3,
		Description:  "description of event 2",
		UserID:       "2",
		Notification: time.Now().AddDate(0, -2, 0),
	}

	ctx := context.Background() // TODO HELP
	_, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)

	eventToUpdate2 := EventToUpdate{
		Title:       "new event2",
		Description: "new description of event 2",
	}

	err = store.Update(id2, eventToUpdate2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)

	val2, ok := store.Events[id2]
	require.True(t, ok)
	require.Equal(t, val2.Title, "new event2")
	require.Equal(t, val2.Description, "new description of event 2")
}

func TestStorageDelete(t *testing.T) {
	store := New()
	eventToCreate1 := EventToCreate{
		DateStart:    time.Now(),
		Title:        "event1",
		DateEnd:      3,
		Description:  "description of event 1",
		UserID:       "1",
		Notification: time.Now().AddDate(0, -1, 0),
	}
	eventToCreate2 := EventToCreate{
		DateStart:    time.Now(),
		Title:        "event2",
		DateEnd:      3,
		Description:  "description of event 2",
		UserID:       "2",
		Notification: time.Now().AddDate(0, -2, 0),
	}
	ctx := context.Background() // TODO HELP
	id1, err := store.Add(ctx, eventToCreate1)
	require.NoError(t, err)
	id2, err := store.Add(ctx, eventToCreate2)
	require.NoError(t, err)

	require.Equal(t, len(store.Events), 2)
	_, ok := store.Events[id1]
	require.True(t, ok)

	store.Delete(id1)
	require.Equal(t, len(store.Events), 1)
	_, ok = store.Events[id2]
	require.True(t, ok)
}

func TestStorageGetEventListing(t *testing.T) {
	store := New()
	date1 := time.Date(2025, time.September, 4, 10, 0o0, 0o0, 0, time.Local)
	date2 := time.Date(2025, time.September, 5, 10, 0o0, 0o0, 0, time.Local)

	eventToCreate1 := EventToCreate{
		DateStart:    date1,
		Title:        "event1",
		DateEnd:      1,
		Description:  "description of event 1",
		UserID:       "1",
		Notification: time.Now().AddDate(0, -1, 0),
	}
	eventToCreate2 := EventToCreate{
		DateStart:    date2,
		Title:        "event2",
		DateEnd:      2,
		Description:  "description of event 2",
		UserID:       "1",
		Notification: time.Now().AddDate(0, -2, 0),
	}
	eventToCreate3 := EventToCreate{
		DateStart:    time.Date(2025, time.September, 20, 10, 0o0, 0o0, 0, time.Local),
		Title:        "event3",
		DateEnd:      3,
		Description:  "description of event 3",
		UserID:       "1",
		Notification: time.Now().AddDate(0, -2, 0),
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
	require.Equal(t, res1[0], store.Events[id2])

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
