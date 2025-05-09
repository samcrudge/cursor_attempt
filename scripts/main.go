package main

import (
	"flag"
	"fmt"
	"os"

	"coding_challenge/scripts/loadtest"
	"coding_challenge/scripts/testclient"
)

func main() {
	// Parse command line flags
	testType := flag.String("type", "client", "Type of test to run: 'client' (default) or 'load'")
	flag.Parse()

	switch *testType {
	case "client":
		fmt.Println("Running test client...")
		testclient.RunTest()
	case "load":
		fmt.Println("Running load test...")
		loadtest.RunTest()
	default:
		fmt.Printf("Unknown test type: %s\n", *testType)
		fmt.Println("Available types: 'client', 'load'")
		os.Exit(1)
	}
}

// Event represents an incoming event
type Event struct {
	ID        string `json:"id"`
	Timestamp int64  `json:"timestamp"`
	Payload   string `json:"payload"`
}
