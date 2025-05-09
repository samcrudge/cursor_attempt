package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"coding_challenge/internal/models"
)

// Server represents the HTTP API server
type Server struct {
	server     *http.Server
	eventStore *models.EventStore
	logger     *log.Logger
}

// NewServer creates a new API server
func NewServer(addr string, eventStore *models.EventStore, logger *log.Logger) *Server {
	router := mux.NewRouter()
	server := &Server{
		server: &http.Server{
			Addr:         addr,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		},
		eventStore: eventStore,
		logger:     logger,
	}

	// Set up routes
	router.HandleFunc("/events", server.handlePostEvent).Methods(http.MethodPost)
	router.HandleFunc("/events", server.handleGetEvents).Methods(http.MethodGet)
	router.HandleFunc("/events/{id}", server.handleGetEvent).Methods(http.MethodGet)
	router.HandleFunc("/health", server.handleHealth).Methods(http.MethodGet)

	return server
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Printf("Starting HTTP server on %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Println("Shutting down HTTP server...")
	return s.server.Shutdown(ctx)
}

// handlePostEvent processes POST requests to create a new event
func (s *Server) handlePostEvent(w http.ResponseWriter, r *http.Request) {
	var event models.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := models.ValidateEvent(&event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.eventStore.Add(&event); err != nil {
		if err == models.ErrDuplicateEventID {
			http.Error(w, "Event with this ID already exists", http.StatusConflict)
		} else {
			http.Error(w, "Failed to store event", http.StatusInternalServerError)
		}
		return
	}

	s.logger.Printf("Received event: %s", event.ID)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": event.ID})
}

// handleGetEvents returns all events
func (s *Server) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	events := s.eventStore.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

// handleGetEvent returns a specific event by ID
func (s *Server) handleGetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	event, err := s.eventStore.Get(id)
	if err != nil {
		if err == models.ErrEventNotFound {
			http.Error(w, "Event not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to retrieve event", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(event)
}

// handleHealth provides a basic health check
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
