package main

import (
  "fmt"
  "net/http"
  "testing"
  "time"

  "github.com/mehXX/awaitility"
)

func startServer() {
  time.Sleep(3 * time.Second)
  http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, world!")
  })
  http.ListenAndServe(":8080", nil)
}

func TestServerRunning(t *testing.T) {
  go startServer()

  err := awaitility.Await(100*time.Millisecond, 5000*time.Millisecond, func() bool {
    resp, err := http.Get("http://localhost:8080")
    if err != nil {
      return false
    }
    defer func() {
      err = resp.Body.Close()
      if err != nil {
        t.Logf("Failed to close response body: %v", err)
      }
    }()
    return resp.StatusCode == http.StatusOK
  })

  if err != nil {
    t.Errorf("Unexpected error during await: %s", err)
  }

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
