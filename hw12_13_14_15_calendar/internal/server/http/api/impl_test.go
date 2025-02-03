package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/app"                          //nolint:depguard
	"github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/logger"                       //nolint:depguard
	memorystorage "github.com/msa16/otus-hw/hw12_13_14_15_calendar/internal/storage/memory" //nolint:depguard
	middleware "github.com/oapi-codegen/nethttp-middleware"                                 //nolint:depguard
	"github.com/oapi-codegen/testutil"                                                      //nolint:depguard
	"github.com/stretchr/testify/assert"                                                    //nolint:depguard
	"github.com/stretchr/testify/require"                                                   //nolint:depguard
)

func doGet(t *testing.T, mux *http.ServeMux, url string) *httptest.ResponseRecorder {
	t.Helper()
	response := testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, mux)
	return response.Recorder
}

func TestCalendar(t *testing.T) {
	testStartTime, _ := time.Parse(time.RFC3339, "2025-01-02T15:00:00Z")
	testStopTime, _ := time.Parse(time.RFC3339, "2025-01-02T16:00:00Z")
	testReminder := "1h0m0s"
	testDescription := "test description"
	allEvents := make(map[string]int64)

	var err error

	// Get the swagger description of our API
	swagger, err := GetSwagger()
	require.NoError(t, err)

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	// Create a new ServeMux for testing.
	m := http.NewServeMux()

	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	opts := StdHTTPServerOptions{
		BaseRouter: m,
		Middlewares: []MiddlewareFunc{
			middleware.OapiRequestValidator(swagger),
		},
	}

	testApp := &app.App{Logger: logger.New("INFO", "/tmp/calendar-test.log"), Storage: memorystorage.New()}
	store := NewAPIServer(testApp)

	HandlerWithOptions(store, opts)

	t.Run("Create events", func(t *testing.T) {
		for i := 0; i < 10; i++ {
			newEvent := NewEvent{
				Title:       "new event " + strconv.Itoa(i),
				Description: &testDescription,
				StartTime:   testStartTime,
				StopTime:    testStopTime,
				UserID:      int64(i),
				Reminder:    &testReminder,
			}

			rr := testutil.NewRequest().Post("/events").WithJsonBody(newEvent).GoWithHTTPHandler(t, m).Recorder
			assert.Equal(t, http.StatusCreated, rr.Code)

			var resultEventID EventID
			err = json.NewDecoder(rr.Body).Decode(&resultEventID)
			assert.NoError(t, err, "error unmarshaling response")
			assert.Equal(t, 36, len(resultEventID.ID), "event id should be 36 chars long")
			allEvents[resultEventID.ID] = int64(i)
		}
	})

	t.Run("Get event", func(t *testing.T) {
		for eventID, userID := range allEvents {
			rr := doGet(t, m, "/events/"+eventID)
			assert.Equal(t, http.StatusOK, rr.Code)

			var resultEvent Event
			err = json.NewDecoder(rr.Body).Decode(&resultEvent)
			assert.NoError(t, err, "error unmarshaling response")
			assert.Equal(t, eventID, resultEvent.ID, "event id should match")
			assert.Equal(t, "new event "+strconv.Itoa(int(userID)), resultEvent.Title, "event title should match")
			assert.Equal(t, testDescription, *resultEvent.Description, "event description should match")
			assert.Equal(t, testStartTime, resultEvent.StartTime, "event start time should match")
			assert.Equal(t, testStopTime, resultEvent.StopTime, "event stop time should match")
			assert.Equal(t, userID, resultEvent.UserID, "event user id should match")
			assert.Equal(t, testReminder, *resultEvent.Reminder, "event reminder should match")
		}
	})

	t.Run("Update event", func(t *testing.T) {
		for eventID, userID := range allEvents {
			newEvent := Event{
				Title:       "event title updated",
				Description: &testDescription,
				StartTime:   testStartTime,
				StopTime:    testStopTime,
				UserID:      userID,
				Reminder:    &testReminder,
				ID:          eventID,
			}

			rr := testutil.NewRequest().Put("/events/"+eventID).WithJsonBody(newEvent).GoWithHTTPHandler(t, m).Recorder
			assert.Equal(t, http.StatusNoContent, rr.Code)
		}
	})

	t.Run("Find events", func(t *testing.T) {
		rr := doGet(t, m, "/events?period=day&startTime="+testStartTime.Format(time.RFC3339))
		assert.Equal(t, http.StatusOK, rr.Code)

		var resultEvents []Event
		err = json.NewDecoder(rr.Body).Decode(&resultEvents)
		assert.NoError(t, err, "error unmarshaling response")
		assert.Equal(t, len(allEvents), len(resultEvents), "event count should be the same as the number of events created")

		for _, event := range resultEvents {
			assert.Equal(t, event.UserID, allEvents[event.ID], "event user id should match")
		}
	})

	t.Run("Delete event", func(t *testing.T) {
		for eventID := range allEvents {
			rr := testutil.NewRequest().Delete("/events/"+eventID).GoWithHTTPHandler(t, m).Recorder
			assert.Equal(t, http.StatusNoContent, rr.Code)
		}
	})
}
