package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
	testReminder := "1h"
	testDescription := "test description"
	var testEventID string

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

	t.Run("Create event", func(t *testing.T) {
		newEvent := NewEvent{
			Title:       "new event",
			Description: &testDescription,
			StartTime:   testStartTime,
			StopTime:    testStopTime,
			UserID:      1,
			Reminder:    &testReminder,
		}

		rr := testutil.NewRequest().Post("/events").WithJsonBody(newEvent).GoWithHTTPHandler(t, m).Recorder
		assert.Equal(t, http.StatusCreated, rr.Code)

		var resultEventID EventID
		err = json.NewDecoder(rr.Body).Decode(&resultEventID)
		assert.NoError(t, err, "error unmarshaling response")
		assert.Equal(t, 36, len(resultEventID.ID), "event id should be 36 chars long")
		testEventID = resultEventID.ID
	})

	t.Run("Get event", func(t *testing.T) {
		rr := doGet(t, m, "/events/"+testEventID)
		assert.Equal(t, http.StatusOK, rr.Code)

		var resultEvent Event
		err = json.NewDecoder(rr.Body).Decode(&resultEvent)
		assert.NoError(t, err, "error unmarshaling response")
	})

	t.Run("Update event", func(t *testing.T) {
		newEvent := Event{
			Title:       "event title",
			Description: &testDescription,
			StartTime:   testStartTime,
			StopTime:    testStopTime,
			UserID:      1,
			Reminder:    &testReminder,
			ID:          testEventID,
		}

		rr := testutil.NewRequest().Put("/events/"+testEventID).WithJsonBody(newEvent).GoWithHTTPHandler(t, m).Recorder
		assert.Equal(t, http.StatusNoContent, rr.Code)
	})
}
