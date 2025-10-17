package integrationtests

import (
	"fmt"
	"net/http"
	"os"
	"slices"
	"testing"
	"time"

	"github.com/c2fo/testify/require"
	"resty.dev/v3"
)

func TestStart(t *testing.T) {
	x := true
	require.True(t, x)
}

func TestHello(t *testing.T) {
	client := resty.New()
	defer client.Close()

	res, err := client.R().Get("http://calendar:8081/")
	require.NoError(t, err)
	require.Equal(t, res.Bytes(), "Hello!")
	require.Equal(t, http.StatusOK, res.StatusCode())
}

func TestGetInexistentEvent(t *testing.T) {
	client := resty.New()
	defer client.Close()

	res, err := client.R().Get("http://calendar:8081/user/1/event/9999999")
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode())
}

func TestPutAndGetEvent(t *testing.T) {
	client := resty.New()
	defer client.Close()

	var createdID map[string]string

	res, err := client.R().
		SetBody(
			map[string]any{
				"title":     "test1",
				"dateStart": time.Now().Add(1 * time.Hour),
				"dateEnd":   time.Now().Add(2 * time.Hour),
			},
		).
		SetResult(&createdID).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	res, err = client.R().Get("http://calendar:8081/user/1/event/" + createdID["id"])
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())

	// clean up
	res, err = client.R().
		SetError(nil).
		Delete("http://calendar:8081/user/1/event/" + createdID["id"])
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())

	res, err = client.R().
		Get("http://calendar:8081/user/1/event/" + createdID["id"])
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, res.StatusCode())
}

func TestPutAndGetListOfEventsForDay(t *testing.T) {
	client := resty.New()
	defer client.Close()

	var createdID1 map[string]string
	var createdID2 map[string]string
	var createdID3 map[string]string

	// -------------------------------------------------------------------------------
	res, err := client.R().
		SetBody(
			map[string]any{
				"title":     "test_day_1",
				"dateStart": time.Now().Add(2 * time.Hour),
				"dateEnd":   time.Now().Add(3 * time.Hour),
			},
		).
		SetResult(&createdID1).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	res, err = client.R().
		SetBody(
			map[string]any{
				"title":     "test_day_2",
				"dateStart": time.Now().Add(2 * time.Hour),
				"dateEnd":   time.Now().Add(3 * time.Hour),
			},
		).
		SetResult(&createdID2).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	res, err = client.R().
		SetBody(
			map[string]any{
				"title":     "test_day_3",
				"dateStart": time.Now().Add(2 * time.Hour),
				"dateEnd":   time.Now().Add(3 * time.Hour),
			},
		).
		SetResult(&createdID3).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	createdIDs := []string{createdID1["id"], createdID2["id"], createdID3["id"]}

	// -------------------------------------------------------------------------------

	res, err = client.R().
		SetQueryParams(map[string]string{
			"period": "day",
		}).
		SetResult(&[]map[string]any{}).
		Get("http://calendar:8081/user/1/events/")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	resultBody := res.Result().(*[]map[string]any)
	require.Len(t, *resultBody, 3)

	resultEventsIDs := []string{}

	for _, event := range *resultBody {
		id, ok := event["ID"].(string)
		require.True(t, ok)
		resultEventsIDs = append(resultEventsIDs, id)
	}

	ok := slices.Equal(resultEventsIDs, createdIDs)
	require.True(t, ok)

	// -------------------------------------------------------------------------------

	// clean up

	for _, id := range createdIDs {
		res, err = client.R().
			SetError(nil).
			Delete("http://calendar:8081/user/1/event/" + id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())

		res, err = client.R().
			Get("http://calendar:8081/user/1/event/" + id)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode())
	}
}

func TestPutAndGetListOfEventsForDayAndCheckNotification(t *testing.T) {
	client := resty.New()
	defer client.Close()

	var createdID1 map[string]string
	var createdID2 map[string]string
	var createdID3 map[string]string

	// -------------------------------------------------------------------------------
	res, err := client.R().
		SetBody(
			map[string]any{
				"title":     "test_day_21",
				"dateStart": time.Now().Add(1 * time.Hour),
				"dateEnd":   time.Now().Add(2 * time.Hour),
			},
		).
		SetResult(&createdID1).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	res, err = client.R().
		SetBody(
			map[string]any{
				"title":     "test_day_22",
				"dateStart": time.Now().Add(1 * time.Hour),
				"dateEnd":   time.Now().Add(2 * time.Hour),
			},
		).
		SetResult(&createdID2).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	res, err = client.R().
		SetBody(
			map[string]any{
				"title":     "test_day_23",
				"dateStart": time.Now().Add(1 * time.Hour),
				"dateEnd":   time.Now().Add(2 * time.Hour),
			},
		).
		SetResult(&createdID3).
		SetError(nil).
		Put("http://calendar:8081/user/1/event/")

	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, res.StatusCode())

	createdIDs := []string{createdID1["id"], createdID2["id"], createdID3["id"]}

	// -------------------------------------------------------------------------------

	res, err = client.R().
		SetQueryParams(map[string]string{
			"period": "day",
		}).
		SetResult(&[]map[string]any{}).
		Get("http://calendar:8081/user/1/events/")
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode())
	resultBody := res.Result().(*[]map[string]any)

	resultEventsIDs := []string{}

	for _, event := range *resultBody {
		id, ok := event["ID"].(string)
		require.True(t, ok)
		resultEventsIDs = append(resultEventsIDs, id)
	}

	ok := slices.Equal(resultEventsIDs, createdIDs)
	require.True(t, ok)

	// -------------------------------------------------------------------------------
	// check that events are notified

	time.Sleep(7 * time.Second)

	for _, id := range createdIDs {
		res, err := client.R().SetResult(&map[string]any{}).
			Get("http://calendar:8081/user/1/event/" + id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())

		resultBody := res.Result().(*map[string]any)
		val, ok := (*resultBody)["notified"]
		require.True(t, ok)
		require.Equal(t, true, val)
	}

	s, err := os.ReadFile("/var/log/sender.log")
	require.NoError(t, err)
	senderLogs := string(s)

	for _, id := range createdIDs {
		require.Contains(t, senderLogs, fmt.Sprintf("\\\"id\\\":\\\"%s\\\"", id))
	}

	// -------------------------------------------------------------------------------

	// clean up

	for _, id := range createdIDs {
		res, err = client.R().
			SetError(nil).
			Delete("http://calendar:8081/user/1/event/" + id)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode())

		res, err = client.R().
			Get("http://calendar:8081/user/1/event/" + id)
		require.NoError(t, err)
		require.Equal(t, http.StatusNotFound, res.StatusCode())
	}
}
