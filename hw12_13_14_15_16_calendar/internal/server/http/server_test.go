package internalhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/mocks"
	"github.com/adettelle/hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/c2fo/testify/require"
	"github.com/golang/mock/gomock"
	"go.uber.org/zap"
)

const layout = "2006-01-02 15:04:05"

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

	baseTime, err := time.ParseInLocation(layout, "2025-09-19 13:01:45", time.UTC)
	require.NoError(t, err)

	createdAt := baseTime
	start := createdAt.AddDate(0, 0, 0)
	end := createdAt.AddDate(0, 0, 2)
	notification := createdAt.AddDate(0, 0, 1)

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")
	request.SetPathValue("id", "1")

	response := httptest.NewRecorder()

	expectedEvent := storage.Event{
		ID:           "1",
		Title:        "title1",
		CreatedAt:    createdAt,
		Start:        start,
		End:          end,
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
	require.True(t, expectedEvent.Start.Equal(actual.Start))
	require.True(t, expectedEvent.End.Equal(actual.End))
	require.Equal(t, expectedEvent.Description, actual.Description)
	require.Equal(t, expectedEvent.UserID, actual.UserID)
	require.True(t, expectedEvent.Notification.Equal(actual.Notification))
}

func TestAddEventByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorager(ctrl)

	eh := &EventHandlers{
		Storager: mockStorage,
	}

	m := eh.Storager.(*mocks.MockStorager)

	userID := "1"
	reqURL := "/user/1/event/"

	baseTime, err := time.ParseInLocation(layout, "2025-09-19 13:01:45", time.UTC)
	require.NoError(t, err)

	createdAt := baseTime
	start := createdAt.AddDate(0, 0, 0)
	end := createdAt.AddDate(0, 0, 2)
	notification := createdAt.AddDate(0, 0, 1)

	event := storage.EventCreateDTO{
		Title:        "title1",
		Start:        start,
		End:          end,
		Description:  "description1",
		Notification: notification,
	}

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)
	reqBody := string(jsonData)
	wantHTTPStatus := 201

	request, err := http.NewRequestWithContext(
		context.Background(), http.MethodPut, reqURL, strings.NewReader(reqBody),
	)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")

	response := httptest.NewRecorder()

	expectedID := "1"
	m.EXPECT().AddEventByID(context.Background(), event, userID).Return(expectedID, nil) // expectedID

	eh.AddEvent(response, request)

	require.Equal(t, wantHTTPStatus, response.Code)
	require.Contains(t, response.Header().Get("Content-Type"), "application/json")

	res := make(map[string]any)
	err = json.NewDecoder(response.Body).Decode(&res)
	require.NoError(t, err)
	require.Equal(t, expectedID, res["id"])
}

func TestUpdateEventByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorager(ctrl)

	eh := &EventHandlers{
		Storager: mockStorage,
	}

	m := eh.Storager.(*mocks.MockStorager)

	userID := "1"
	eventID := "1"
	reqURL := "/update/user/1/event/1"

	baseTime, err := time.ParseInLocation(layout, "2025-09-19 13:01:45", time.UTC)
	require.NoError(t, err)

	date := baseTime.AddDate(0, 0, 0)
	dateEnd := baseTime.AddDate(0, 0, 1)

	title := "new_title"
	description := "new_description"

	event := storage.EventUpdateDTO{
		Title:       &title,
		Description: &description,
		Start:       &date,
		End:         &dateEnd,
	}

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)
	reqBody := string(jsonData)
	wantHTTPStatus := 202

	request, err := http.NewRequestWithContext(
		context.Background(), http.MethodPost, reqURL, strings.NewReader(reqBody),
	)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")
	request.SetPathValue("id", "1")

	response := httptest.NewRecorder()

	// expectedID := "1"
	m.EXPECT().UpdateEventByID(context.Background(), eventID, event, userID).Return(nil) // expectedID

	eh.UpdateEventeByID(response, request)

	require.Equal(t, wantHTTPStatus, response.Code)
}

func TestUpdateEventWithWrongEndDate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorager(ctrl)
	l, err := zap.NewDevelopment()
	require.NoError(t, err)
	eh := &EventHandlers{
		Storager: mockStorage,
		Logg:     l,
	}

	m := eh.Storager.(*mocks.MockStorager)

	userID := "1"
	eventID := "1"
	reqURL := "/update/user/1/event/1"

	baseTime, err := time.ParseInLocation(layout, "2025-09-19 13:01:45", time.UTC)
	require.NoError(t, err)

	date := baseTime.AddDate(0, 0, 0)
	dateEnd := baseTime.AddDate(0, 0, -1)

	title := "new_title"
	description := "new_description"

	event := storage.EventUpdateDTO{
		Title:       &title,
		Description: &description,
		Start:       &date,
		End:         &dateEnd,
	}

	jsonData, err := json.Marshal(event)
	require.NoError(t, err)
	reqBody := string(jsonData)
	wantHTTPStatus := 400

	request, err := http.NewRequestWithContext(
		context.Background(), http.MethodPost, reqURL, strings.NewReader(reqBody),
	)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")
	request.SetPathValue("id", "1")

	response := httptest.NewRecorder()

	// expectedID := "1"
	m.EXPECT().UpdateEventByID(context.Background(), eventID, event, userID).Times(0) // expectedID

	eh.UpdateEventeByID(response, request)

	require.Equal(t, wantHTTPStatus, response.Code)
}

func TestDeleteEventByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorager(ctrl)

	eh := &EventHandlers{
		Storager: mockStorage,
	}

	m := eh.Storager.(*mocks.MockStorager)

	eventID := "1"
	reqURL := "/user/1/event/1"

	wantHTTPStatus := 200

	request, err := http.NewRequestWithContext(
		context.Background(), http.MethodDelete, reqURL, nil,
	)
	require.NoError(t, err)
	request.SetPathValue("id", "1")

	response := httptest.NewRecorder()

	m.EXPECT().DeleteEventByID(context.Background(), eventID).Return(nil)

	eh.DeleteEventByID(response, request)

	require.Equal(t, wantHTTPStatus, response.Code)
}

func TestGetEventListByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockStorager(ctrl)

	eh := &EventHandlers{
		Storager: mockStorage,
	}

	m := eh.Storager.(*mocks.MockStorager)
	userID := "1"

	reqURL := "/user/1/events/?date=2025-09-19&period=day"

	baseTime, err := time.ParseInLocation(layout, "2025-09-19 00:00:00", time.UTC)
	require.NoError(t, err)

	createdAt := baseTime
	start := createdAt.AddDate(0, 0, 0)
	end := createdAt.AddDate(0, 0, 2)
	notification := createdAt.AddDate(0, 0, 1)

	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, reqURL, nil)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")
	parsedTime := baseTime
	period := "day"

	response := httptest.NewRecorder()

	expectedEvent1 := storage.Event{
		ID:           "1",
		Title:        "title1",
		CreatedAt:    createdAt,
		Start:        start,
		End:          end,
		Description:  "description1",
		UserID:       userID,
		Notification: notification,
	}
	expectedEvent2 := storage.Event{
		ID:           "2",
		Title:        "title2",
		CreatedAt:    createdAt.Add(time.Hour),
		Start:        start.Add(5 * time.Hour),
		End:          end,
		Description:  "description2",
		UserID:       userID,
		Notification: notification,
	}
	expectedEvents := []storage.Event{expectedEvent1, expectedEvent2}

	m.EXPECT().GetEventListingByUserID(userID, parsedTime, period).Return(expectedEvents, nil)

	eh.GetEventListingByUserID(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Contains(t, response.Header().Get("Content-Type"), "application/json")

	var actual []storage.Event
	err = json.NewDecoder(response.Body).Decode(&actual)
	require.NoError(t, err)
	require.Len(t, actual, 2)
	// require.Equal(t, expectedEvent.ID, actual.ID)
	// require.Equal(t, expectedEvent.Title, actual.Title)
	// require.True(t, expectedEvent.CreatedAt.Equal(actual.CreatedAt))
	// require.True(t, expectedEvent.Date.Equal(actual.Date))
	// require.True(t, expectedEvent.Duration.Equal(actual.Duration))
	// require.Equal(t, expectedEvent.Description, actual.Description)
	// require.Equal(t, expectedEvent.UserID, actual.UserID)
	// require.True(t, expectedEvent.Notification.Equal(actual.Notification))
}
