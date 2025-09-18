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

	t1 := time.Now()
	// t2 := t1.Add(-1 * time.Duration.Day)

	d1 := "2025-09-17T20:41:10.411997Z"
	layout := time.RFC3339
	createdAt, err := time.Parse(layout, d1)
	require.NoError(t, err)

	d2 := "2025-09-18T20:41:10.411997Z"
	date, err := time.Parse(layout, d2)
	require.NoError(t, err)

	d3 := "2025-09-19T20:41:10.411997Z"
	duration, err := time.Parse(layout, d3)
	require.NoError(t, err)

	d4 := "2025-09-17T20:41:10.411997Z"
	notification, err := time.Parse(layout, d4)
	require.NoError(t, err)

	m.EXPECT().GetEventByID(eventID, userID).Return(storage.Event{
		ID:           "1",
		Title:        "title1",
		CreatedAt:    createdAt,
		Date:         date,
		Duration:     duration,
		Description:  "description1",
		UserID:       "1",
		Notification: notification,
	}, nil)

	request, err := http.NewRequest(http.MethodGet, reqURL, nil)
	require.NoError(t, err)
	request.SetPathValue("userid", "1")
	request.SetPathValue("id", "1")

	response := httptest.NewRecorder()

	expectedRes := map[string]any{
		"ID":    "1",
		"Tilte": "tile1",
		"Date":  t1.String(),
	}

	actualRes := map[string]any{}
	// x, _ := json.Marshal(expectedRes)

	eh.GetEventByID(response, request)
	result := response.Result()
	json.Unmarshal(response.Body.Bytes(), &actualRes)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-type"))
	// x := `{"ID":"1","Title":"title1","CreatedAt":"2025-09-17T20:41:10.411997Z","Date":"2025-09-18T20:41:10.411997Z","Duration":"2025-09-19T20:41:10.411997Z","Description":"description1","UserID":"1","Notification":"2025-09-17T20:41:10.411997Z"}`
	//require.Equal(t, x, response.Body.String())
	require.Equal(t, expectedRes, actualRes)
	defer result.Body.Close()
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
