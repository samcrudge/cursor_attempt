package models

import (
	"sync"
)

// EventStore provides thread-safe storage and retrieval of events
type EventStore struct {
	mu     sync.RWMutex
	events map[string]*Event
	// Channel for new events
	eventCh chan *Event
}

// NewEventStore creates a new event store with a buffer for event channel
func NewEventStore(bufferSize int) *EventStore {
	return &EventStore{
		events:  make(map[string]*Event),
		eventCh: make(chan *Event, bufferSize),
	}
}

// Add stores an event in the in-memory store
func (s *EventStore) Add(event *Event) error {
	if event == nil {
		return ErrMissingID
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[event.ID]; exists {
		return ErrDuplicateEventID
	}

	s.events[event.ID] = event
	// Send to channel (non-blocking)
	select {
	case s.eventCh <- event:
		// Event sent successfully
	default:
		// Channel is full, just continue
	}

	return nil
}

// Get retrieves an event by ID
func (s *EventStore) Get(id string) (*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return nil, ErrEventNotFound
	}
	return event, nil
}

// GetAll returns all events
func (s *EventStore) GetAll() []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}
	return events
}

// EventChannel returns the channel that emits new events
func (s *EventStore) EventChannel() <-chan *Event {
	return s.eventCh
}

// Close shuts down the event store and its channels
func (s *EventStore) Close() {
	close(s.eventCh)
}
