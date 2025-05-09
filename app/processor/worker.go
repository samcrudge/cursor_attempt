package processor

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/google/uuid"

	"coding_challenge/internal/models"
)

// Worker represents a background processor for events
type Worker struct {
	id         string
	eventStore *models.EventStore
	logger     *log.Logger
}

// NewWorker creates a new background worker
func NewWorker(eventStore *models.EventStore, logger *log.Logger) *Worker {
	return &Worker{
		id:         uuid.New().String()[:8], // short worker ID
		eventStore: eventStore,
		logger:     logger,
	}
}

// Start begins the worker processing loop
func (w *Worker) Start(ctx context.Context) {
	w.logger.Printf("Starting worker %s", w.id)

	eventCh := w.eventStore.EventChannel()

	for {
		select {
		case <-ctx.Done():
			w.logger.Printf("Worker %s shutting down...", w.id)
			return
		case event, ok := <-eventCh:
			if !ok {
				w.logger.Printf("Event channel closed, worker %s shutting down", w.id)
				return
			}
			w.processEvent(event)
		}
	}
}

// processEvent transforms and publishes an event
func (w *Worker) processEvent(event *models.Event) {
	// Transform the event (uppercase the payload)
	transformedEvent := &models.TransformedEvent{
		ID:           event.ID,
		OriginalTime: event.Timestamp,
		ProcessedAt:  time.Now(),
		Payload:      strings.ToUpper(event.Payload), // Simple transformation
		ProcessorID:  w.id,
	}

	// "Publish" the transformed event (in this case, just log it)
	w.publishEvent(transformedEvent)
}

// publishEvent outputs the transformed event
func (w *Worker) publishEvent(event *models.TransformedEvent) {
	w.logger.Printf("[PUBLISHED] Worker %s processed event %s: %s",
		w.id,
		event.ID,
		event.Payload)
}
