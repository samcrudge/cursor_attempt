package loadtest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Event represents an incoming event to be processed
type Event struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Payload   string `json:"payload"`
}

// RunTest performs a load test on the event processing application
func RunTest() {
	// Configure test parameters
	apiURL := "http://localhost:8080/events"
	concurrentClients := 10
	eventsPerClient := 100
	requestTimeout := 2 * time.Second

	logger := log.New(os.Stdout, "[LOAD-TEST] ", log.LstdFlags)
	logger.Printf("Starting load test with %d concurrent clients, each sending %d events\n",
		concurrentClients, eventsPerClient)

	// Start benchmark
	startTime := time.Now()

	// Create wait group to track clients
	var wg sync.WaitGroup
	wg.Add(concurrentClients)

	// Success/failure counters
	var (
		successCount, errorCount int
		counterMutex             sync.Mutex
		rateLimiter              = make(chan struct{}, 20) // Limit to 20 concurrent requests
	)

	// Prepare payloads with varying sizes
	payloads := generatePayloads()

	// Launch clients
	for i := 0; i < concurrentClients; i++ {
		clientID := fmt.Sprintf("client-%d", i)

		go func(id string) {
			defer wg.Done()

			// Create HTTP client with timeout
			client := &http.Client{
				Timeout: requestTimeout,
			}

			for j := 0; j < eventsPerClient; j++ {
				// Rate limiting
				rateLimiter <- struct{}{}
				defer func() { <-rateLimiter }()

				// Create event with random payload
				event := Event{
					ID:        uuid.New().String(),
					Timestamp: time.Now().Unix(),
					Payload:   payloads[j%len(payloads)],
				}

				// Marshal to JSON
				eventJSON, err := json.Marshal(event)
				if err != nil {
					logger.Printf("[%s] Error marshaling event: %v", id, err)
					counterMutex.Lock()
					errorCount++
					counterMutex.Unlock()
					continue
				}

				// Send request
				resp, err := client.Post(apiURL, "application/json", bytes.NewBuffer(eventJSON))
				if err != nil {
					logger.Printf("[%s] Request error: %v", id, err)
					counterMutex.Lock()
					errorCount++
					counterMutex.Unlock()
					continue
				}

				// Check response
				if resp.StatusCode == http.StatusCreated {
					counterMutex.Lock()
					successCount++
					counterMutex.Unlock()
				} else {
					logger.Printf("[%s] Failed response: %d", id, resp.StatusCode)
					counterMutex.Lock()
					errorCount++
					counterMutex.Unlock()
				}
				resp.Body.Close()
			}
		}(clientID)
	}

	// Wait for all clients to finish
	wg.Wait()
	elapsed := time.Since(startTime)

	// Print results
	totalRequests := successCount + errorCount
	rps := float64(totalRequests) / elapsed.Seconds()

	logger.Printf("\n=== Load Test Results ===")
	logger.Printf("Duration: %.2f seconds", elapsed.Seconds())
	logger.Printf("Total Requests: %d", totalRequests)
	logger.Printf("Successful Requests: %d (%.2f%%)",
		successCount, float64(successCount)*100/float64(totalRequests))
	logger.Printf("Failed Requests: %d (%.2f%%)",
		errorCount, float64(errorCount)*100/float64(totalRequests))
	logger.Printf("Requests per Second: %.2f", rps)
	logger.Printf("=========================")

	// Check server health
	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		logger.Printf("Error checking server health: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		logger.Printf("Server health: OK")
	} else {
		logger.Printf("Server health: NOT OK (status: %d)", resp.StatusCode)
	}
}

// generatePayloads creates test data with varying sizes
func generatePayloads() []string {
	return []string{
		"small payload",
		"This is a medium sized payload with some more text to process",
		"This is a larger payload that contains more data to process and would take slightly more CPU time to transform. It simulates a more realistic event payload that might contain various fields and information that needs processing.",
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vestibulum eget ligula eu lectus finibus condimentum. Mauris vestibulum iaculis eros, eu venenatis enim vehicula at. Integer nec orci in mi vulputate dictum quis in nisl. Nulla facilisi. Donec consequat dignissim augue, ultricies dapibus urna facilisis et. Sed pulvinar, tellus in rhoncus pharetra, neque massa facilisis eros, at rutrum tortor nunc in urna. Proin suscipit augue at orci feugiat lobortis.",
	}
}
