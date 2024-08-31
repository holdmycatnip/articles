package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"
)

// as an example we are starting a server with some delay, but it could be any other async operation that would take some time
// like publishing message to message queue.
func startServer() {
	rand.Seed(time.Now().UnixNano())
	sleepDuration := time.Duration(rand.Intn(5)+1) * time.Second
	time.Sleep(sleepDuration)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world!")
	})
	http.ListenAndServe(":8080", nil)
}

func TestServerRunning(t *testing.T) {
	go startServer()

	time.Sleep(2 * time.Second)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			t.Logf("Failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.StatusCode)
	}
}
