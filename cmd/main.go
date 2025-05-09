package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"coding_challenge/app/api"
	"coding_challenge/app/processor"
	"coding_challenge/internal/models"
)

const (
	workerCount     = 3
	eventBufferSize = 100
	serverAddress   = ":8081"
)

func main() {
	// Set up logger
	logger := log.New(os.Stdout, "[EVENT-PROCESSOR] ", log.LstdFlags)
	logger.Println("Starting event processor application...")

	// Create a context that will be canceled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Start components with waitgroup to track active components
	var wg sync.WaitGroup

	// Initialize event store
	eventStore := models.NewEventStore(eventBufferSize)

	// Start API server
	apiServer := api.NewServer(serverAddress, eventStore, logger)
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := apiServer.Start(); err != nil && err != http.ErrServerClosed {
			logger.Printf("HTTP server error: %v", err)
		}
	}()

	// Start worker(s)
	for i := 0; i < workerCount; i++ {
		worker := processor.NewWorker(eventStore, logger)
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker.Start(ctx)
		}()
	}

	// Wait for shutdown signal
	sig := <-sigCh
	logger.Printf("Received signal %v, initiating graceful shutdown...", sig)

	// Cancel context to notify all components to shut down
	cancel()

	// Set a timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown the HTTP server
	if err := apiServer.Stop(shutdownCtx); err != nil {
		logger.Printf("Error during server shutdown: %v", err)
	}

	// Close the event store
	eventStore.Close()

	// Wait for all components to shut down or timeout
	shutdownCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(shutdownCh)
	}()

	select {
	case <-shutdownCh:
		logger.Println("All components shut down successfully")
	case <-shutdownCtx.Done():
		logger.Println("Shutdown timed out, forcing exit")
	}

	logger.Println("Application stopped")
}
