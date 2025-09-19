package internalhttp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/mocks"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/c2fo/testify/require"
	"github.com/golang/mock/gomock"
)

func TestGetEventByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorager(ctrl)

	eh := &EventHandlers{
		Storager: mockStorage,
	}

	m := eh.Storager.(*mocks.MockStorager)
	eventID := "1"
	userID := "1"

	reqURL := "/user/1/event/1"

	layout := "2006-01-02 15:04:05"
	baseTime, err := time.ParseInLocation(layout, "2025-09-19 13:01:45", time.UTC) // или time.Local, либо time.FixedZone(...)
	require.NoError(t, err)

	createdAt := baseTime
	date := createdAt.AddDate(0, 0, 0)
	duration := createdAt.AddDate(0, 0, 2)
	notification := createdAt.AddDate(0, 0, 1)

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")
	request.SetPathValue("id", "1")

	response := httptest.NewRecorder()

	expectedEvent := storage.Event{
		ID:           "1",
		Title:        "title1",
		CreatedAt:    createdAt,
		Date:         date,
		Duration:     duration,
		Description:  "description1",
		UserID:       "1",
		Notification: notification,
	}

	m.EXPECT().GetEventByID(eventID, userID).Return(expectedEvent, nil)

	eh.GetEventByID(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Header().Get("Content-Type"), "application/json")

	var actual storage.Event
	err = json.NewDecoder(response.Body).Decode(&actual)
	require.NoError(t, err)

	require.Equal(t, expectedEvent.ID, actual.ID)
	require.Equal(t, expectedEvent.Title, actual.Title)
	require.True(t, expectedEvent.CreatedAt.Equal(actual.CreatedAt))
	require.True(t, expectedEvent.Date.Equal(actual.Date))
	require.True(t, expectedEvent.Duration.Equal(actual.Duration))
	require.Equal(t, expectedEvent.Description, actual.Description)
	require.Equal(t, expectedEvent.UserID, actual.UserID)
	require.True(t, expectedEvent.Notification.Equal(actual.Notification))
}

// func TestAddEventByID(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockStorage := mocks.NewMockStorager(ctrl)

// 	eh := &EventHandlers{
// 		Storager: mockStorage,
// 	}

// 	m := eh.Storager.(*mocks.MockStorager)
// 	userID := "1"

// 	reqURL := "/user/1/event/"

// }

/*
func TestCheckConnectionToDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mDB := mocks.NewMockDBConnector(ctrl)

	eh := &EventHandlers{}

	reqURL := "/"

	mDB.EXPECT().Connect().Return(nil, nil)

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	require.NoError(t, err)

	response := httptest.NewRecorder()

	eh.CheckConnectionToDB(response, request)
	require.Equal(t, http.StatusOK, response.Code)
}
*/
