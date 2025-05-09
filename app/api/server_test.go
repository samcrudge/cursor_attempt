package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"coding_challenge/internal/models"
)

func TestHandlePostEvent(t *testing.T) {
	// Create test event store and server
	eventStore := models.NewEventStore(10)
	logger := log.New(io.Discard, "", 0) // Silent logger for tests
	server := NewServer(":8080", eventStore, logger)

	// Test cases
	testCases := []struct {
		name       string
		eventJSON  string
		wantStatus int
	}{
		{
			name:       "valid event",
			eventJSON:  `{"id":"test-1","timestamp":1625097600,"payload":"test payload"}`,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "missing ID",
			eventJSON:  `{"timestamp":1625097600,"payload":"test payload"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "duplicate ID",
			eventJSON:  `{"id":"test-1","timestamp":1625097600,"payload":"test payload"}`,
			wantStatus: http.StatusConflict,
		},
		{
			name:       "invalid JSON",
			eventJSON:  `{"id":"test-2","timestamp":1625097600,"payload":`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with JSON body
			req := httptest.NewRequest(http.MethodPost, "/events", bytes.NewBufferString(tc.eventJSON))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			rec := httptest.NewRecorder()

			// Handle request
			server.handlePostEvent(rec, req)

			// Check status code
			if rec.Code != tc.wantStatus {
				t.Errorf("Expected status %d, got %d", tc.wantStatus, rec.Code)
			}
		})
	}
}

func TestHandleGetEvents(t *testing.T) {
	// Create test event store and server
	eventStore := models.NewEventStore(10)
	logger := log.New(os.Stdout, "", 0)
	server := NewServer(":8080", eventStore, logger)

	// Add some test events
	events := []*models.Event{
		{ID: "id1", Timestamp: 1625097600, Payload: "payload1"},
		{ID: "id2", Timestamp: 1625097601, Payload: "payload2"},
	}

	for _, e := range events {
		_ = eventStore.Add(e)
	}

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/events", nil)

	// Create response recorder
	rec := httptest.NewRecorder()

	// Handle request
	server.handleGetEvents(rec, req)

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Decode response
	var responseEvents []*models.Event
	if err := json.NewDecoder(rec.Body).Decode(&responseEvents); err != nil {
		t.Errorf("Failed to decode response: %v", err)
	}

	// Check event count
	if len(responseEvents) != len(events) {
		t.Errorf("Expected %d events, got %d", len(events), len(responseEvents))
	}
}

func TestHandleHealth(t *testing.T) {
	// Create test server
	eventStore := models.NewEventStore(10)
	logger := log.New(io.Discard, "", 0)
	server := NewServer(":8080", eventStore, logger)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Create response recorder
	rec := httptest.NewRecorder()

	// Handle request
	server.handleHealth(rec, req)

	// Check status code
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// Check response body contains status
	if !strings.Contains(rec.Body.String(), "status") {
		t.Error("Response body does not contain status field")
	}
}
