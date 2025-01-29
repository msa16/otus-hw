//go:build go1.22

// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/oapi-codegen/runtime"
)

// Error defines model for Error.
type Error struct {
	// Message error message
	Message string `json:"message"`
}

// Event defines model for Event.
type Event struct {
	Description *string   `json:"Description,omitempty"`
	ID          string    `json:"ID"`
	Reminder    *string   `json:"Reminder,omitempty"`
	StartTime   time.Time `json:"StartTime"`
	StopTime    time.Time `json:"StopTime"`
	Title       string    `json:"Title"`
	UserID      int64     `json:"UserID"`
}

// NewEvent defines model for NewEvent.
type NewEvent struct {
	Description *string   `json:"Description,omitempty"`
	Reminder    *string   `json:"Reminder,omitempty"`
	StartTime   time.Time `json:"StartTime"`
	StopTime    time.Time `json:"StopTime"`
	Title       string    `json:"Title"`
	UserID      int64     `json:"UserID"`
}

// FindEventsParams defines parameters for FindEvents.
type FindEventsParams struct {
	// StartTime events start time
	StartTime time.Time `form:"startTime" json:"startTime"`

	// Period period from startTime - day, week, month
	Period *string `form:"period,omitempty" json:"period,omitempty"`
}

// CreateEventJSONRequestBody defines body for CreateEvent for application/json ContentType.
type CreateEventJSONRequestBody = NewEvent

// UpdateEventByIdJSONRequestBody defines body for UpdateEventById for application/json ContentType.
type UpdateEventByIdJSONRequestBody = Event

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all events
	// (GET /events)
	FindEvents(w http.ResponseWriter, r *http.Request, params FindEventsParams)
	// Create new event
	// (POST /events)
	CreateEvent(w http.ResponseWriter, r *http.Request)
	// Delete event by id
	// (DELETE /events/{id})
	DeleteEventById(w http.ResponseWriter, r *http.Request, id string)
	// Get event by id
	// (GET /events/{id})
	FindEventById(w http.ResponseWriter, r *http.Request, id string)
	// Update event by id
	// (PUT /events/{id})
	UpdateEventById(w http.ResponseWriter, r *http.Request, id string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// FindEvents operation middleware
func (siw *ServerInterfaceWrapper) FindEvents(w http.ResponseWriter, r *http.Request) {

	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params FindEventsParams

	// ------------- Required query parameter "startTime" -------------

	if paramValue := r.URL.Query().Get("startTime"); paramValue != "" {

	} else {
		siw.ErrorHandlerFunc(w, r, &RequiredParamError{ParamName: "startTime"})
		return
	}

	err = runtime.BindQueryParameter("form", true, true, "startTime", r.URL.Query(), &params.StartTime)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "startTime", Err: err})
		return
	}

	// ------------- Optional query parameter "period" -------------

	err = runtime.BindQueryParameter("form", true, false, "period", r.URL.Query(), &params.Period)
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "period", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindEvents(w, r, params)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// CreateEvent operation middleware
func (siw *ServerInterfaceWrapper) CreateEvent(w http.ResponseWriter, r *http.Request) {

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.CreateEvent(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// DeleteEventById operation middleware
func (siw *ServerInterfaceWrapper) DeleteEventById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", r.PathValue("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.DeleteEventById(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// FindEventById operation middleware
func (siw *ServerInterfaceWrapper) FindEventById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", r.PathValue("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.FindEventById(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

// UpdateEventById operation middleware
func (siw *ServerInterfaceWrapper) UpdateEventById(w http.ResponseWriter, r *http.Request) {

	var err error

	// ------------- Path parameter "id" -------------
	var id string

	err = runtime.BindStyledParameterWithOptions("simple", "id", r.PathValue("id"), &id, runtime.BindStyledParameterOptions{ParamLocation: runtime.ParamLocationPath, Explode: false, Required: true})
	if err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{ParamName: "id", Err: err})
		return
	}

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.UpdateEventById(w, r, id)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{})
}

// ServeMux is an abstraction of http.ServeMux.
type ServeMux interface {
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type StdHTTPServerOptions struct {
	BaseURL          string
	BaseRouter       ServeMux
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, m ServeMux) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseRouter: m,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, m ServeMux, baseURL string) http.Handler {
	return HandlerWithOptions(si, StdHTTPServerOptions{
		BaseURL:    baseURL,
		BaseRouter: m,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options StdHTTPServerOptions) http.Handler {
	m := options.BaseRouter

	if m == nil {
		m = http.NewServeMux()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	m.HandleFunc("GET "+options.BaseURL+"/events", wrapper.FindEvents)
	m.HandleFunc("POST "+options.BaseURL+"/events", wrapper.CreateEvent)
	m.HandleFunc("DELETE "+options.BaseURL+"/events/{id}", wrapper.DeleteEventById)
	m.HandleFunc("GET "+options.BaseURL+"/events/{id}", wrapper.FindEventById)
	m.HandleFunc("PUT "+options.BaseURL+"/events/{id}", wrapper.UpdateEventById)

	return m
}
