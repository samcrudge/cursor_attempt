package models

import (
	"time"
)

// Event represents an incoming event to be processed
type Event struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Payload   string `json:"payload"`
}

// ValidateEvent checks if an event has all required fields
func ValidateEvent(e *Event) error {
	if e.ID == "" {
		return ErrMissingID
	}
	if e.Timestamp == 0 {
		e.Timestamp = time.Now().Unix()
	}
	return nil
}

// TransformedEvent represents a processed event
type TransformedEvent struct {
	ID           string    `json:"id"`
	OriginalTime int64     `json:"original_time"`
	ProcessedAt  time.Time `json:"processed_at"`
	Payload      string    `json:"payload"`
	ProcessorID  string    `json:"processor_id"`
}

// Common errors
var (
	ErrMissingID        = Error("missing event ID")
	ErrEventNotFound    = Error("event not found")
	ErrDuplicateEventID = Error("duplicate event ID")
)

// Error is a simple string-based error type
type Error string

func (e Error) Error() string {
	return string(e)
}
