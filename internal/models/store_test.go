package models

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestEventStoreAdd(t *testing.T) {
	store := NewEventStore(10)

	event := &Event{
		ID:        "test-id",
		Timestamp: time.Now().Unix(),
		Payload:   "test payload",
	}

	// Test adding a new event
	err := store.Add(event)
	if err != nil {
		t.Errorf("Failed to add event: %v", err)
	}

	// Test adding a duplicate event
	err = store.Add(event)
	if err != ErrDuplicateEventID {
		t.Errorf("Expected duplicate ID error, got: %v", err)
	}

	// Test adding nil event
	err = store.Add(nil)
	if err != ErrMissingID {
		t.Errorf("Expected missing ID error, got: %v", err)
	}
}

func TestEventStoreGet(t *testing.T) {
	store := NewEventStore(10)

	event := &Event{
		ID:        "test-id",
		Timestamp: time.Now().Unix(),
		Payload:   "test payload",
	}

	// Add an event
	_ = store.Add(event)

	// Test getting an existing event
	retrieved, err := store.Get("test-id")
	if err != nil {
		t.Errorf("Failed to get event: %v", err)
	}
	if retrieved.ID != event.ID {
		t.Errorf("Expected ID %s, got %s", event.ID, retrieved.ID)
	}

	// Test getting a non-existent event
	_, err = store.Get("non-existent")
	if err != ErrEventNotFound {
		t.Errorf("Expected event not found error, got: %v", err)
	}
}

func TestEventStoreGetAll(t *testing.T) {
	store := NewEventStore(10)

	// Add multiple events
	events := []*Event{
		{ID: "id1", Timestamp: time.Now().Unix(), Payload: "payload1"},
		{ID: "id2", Timestamp: time.Now().Unix(), Payload: "payload2"},
		{ID: "id3", Timestamp: time.Now().Unix(), Payload: "payload3"},
	}

	for _, e := range events {
		_ = store.Add(e)
	}

	// Get all events and check count
	allEvents := store.GetAll()
	if len(allEvents) != len(events) {
		t.Errorf("Expected %d events, got %d", len(events), len(allEvents))
	}

	// Check if all events are present
	ids := make(map[string]bool)
	for _, e := range allEvents {
		ids[e.ID] = true
	}

	for _, e := range events {
		if !ids[e.ID] {
			t.Errorf("Event %s not found in GetAll results", e.ID)
		}
	}
}

func TestEventStoreConcurrency(t *testing.T) {
	store := NewEventStore(100)
	workers := 10
	eventsPerWorker := 10

	var wg sync.WaitGroup
	wg.Add(workers)

	// Add events concurrently
	for i := 0; i < workers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < eventsPerWorker; j++ {
				id := fmt.Sprintf("worker-%d-event-%d", workerID, j)
				event := &Event{
					ID:        id,
					Timestamp: time.Now().Unix(),
					Payload:   fmt.Sprintf("payload-%d-%d", workerID, j),
				}
				_ = store.Add(event)
			}
		}(i)
	}

	wg.Wait()

	// Check if all events were stored correctly
	allEvents := store.GetAll()
	expectedCount := workers * eventsPerWorker
	if len(allEvents) != expectedCount {
		t.Errorf("Expected %d events, got %d", expectedCount, len(allEvents))
	}
}
