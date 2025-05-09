package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// Event represents an incoming event to be processed
type Event struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Payload   string `json:"payload"`
}

func main() {
	// Configure client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Configure test parameters
	apiURL := "http://localhost:8081/events"
	eventCount := 20
	delayBetweenRequests := 100 * time.Millisecond

	logger := log.New(log.Writer(), "[TEST-CLIENT] ", log.LstdFlags)
	logger.Printf("Starting test client, sending %d events to %s\n", eventCount, apiURL)

	// Generate and send events
	payloads := []string{
		"Hello, world!",
		"Event processing is fun",
		"Cloud native applications",
		"Distributed systems",
		"Concurrent processing",
		"Microservices architecture",
		"Event-driven design",
		"Stream processing",
	}

	for i := 0; i < eventCount; i++ {
		// Create a random event
		event := Event{
			ID:        uuid.New().String(),
			Timestamp: time.Now().Unix(),
			Payload:   payloads[i%len(payloads)],
		}

		// Convert to JSON
		eventJSON, err := json.Marshal(event)
		if err != nil {
			logger.Printf("Error marshaling event: %v", err)
			continue
		}

		// Post to API
		resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(eventJSON))
		if err != nil {
			logger.Printf("Error sending event: %v", err)
			continue
		}

		// Check response
		if resp.StatusCode == http.StatusCreated {
			logger.Printf("Event %s sent successfully", event.ID)
		} else {
			body := make([]byte, 100)
			resp.Body.Read(body)
			logger.Printf("Failed to send event %s: %d - %s", event.ID, resp.StatusCode, string(body))
		}
		resp.Body.Close()

		// Wait before sending next event
		time.Sleep(delayBetweenRequests)
	}

	// Get all events to verify
	logger.Println("Retrieving all events:")
	resp, err := client.Get(apiURL)
	if err != nil {
		logger.Fatalf("Error getting events: %v", err)
	}
	defer resp.Body.Close()

	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		logger.Fatalf("Error decoding events: %v", err)
	}

	logger.Printf("Retrieved %d events", len(events))
	for i, e := range events {
		fmt.Printf("%d. ID: %s, Payload: %s\n", i+1, e.ID, e.Payload)
	}
}
