package node

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandlePut(t *testing.T) {
	// Initialize the node with an empty store
	n := &Node{
		ID:    "localhost:8080",
		Addr:  ":8080",
		Store: NewKeyValueStore(),
	}

	// Prepare request body
	payload := KVPayload{
		Key:   "key1",
		Value: "value1",
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPut, "/put", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	// Record response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(n.HandlePut)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

}

func TestHandleGet(t *testing.T) {
	n := &Node{
		ID:    "localhost:8080",
		Addr:  ":8080",
		Store: NewKeyValueStore(),
	}

	// Put a key-value pair into the store
	n.Store.Put("key1", "value1")

	req, err := http.NewRequest(http.MethodGet, "/get?key=key1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(n.HandleGet)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check response body
	expected := "value1"
	if rr.Body.String() != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, rr.Body.String())
	}
}

func TestHandleDelete(t *testing.T) {
	n := &Node{
		ID:    "localhost:8080",
		Addr:  ":8080",
		Store: NewKeyValueStore(),
	}

	// Put a key-value pair into the store
	n.Store.Put("key1", "value1")

	req, err := http.NewRequest(http.MethodDelete, "/delete?key=key1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record response
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(n.HandleDelete)

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check status code
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	// Check that the key is actually deleted
	_, exists := n.Store.Get("key1")
	if exists {
		t.Error("Expected key 'key1' to be deleted, but it still exists")
	}
}
