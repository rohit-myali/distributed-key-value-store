package node

import (
	"encoding/json"
	"net/http"
)

type Node struct {
	ID    string
	Addr  string
	Store *KeyValueStore
}


type KVPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
// HandlePut stores a key-value pair
func (n *Node) HandlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload KVPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	key, value := payload.Key, payload.Value
	if key == "" || value == "" {
		http.Error(w, "Missing key or value", http.StatusBadRequest)
		return
	}

	n.Store.Put(key, value)
	w.WriteHeader(http.StatusOK)
}

// HandleGet retrieves a value by key
func (n *Node) HandleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}
	value, ok := n.Store.Get(key)
	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.Write([]byte(value))
}

// HandleDelete deletes a key-value pair
func (n *Node) HandleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "Missing key", http.StatusBadRequest)
		return
	}

	ok := n.Store.Delete(key)

	if !ok {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
