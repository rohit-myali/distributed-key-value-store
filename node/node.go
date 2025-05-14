package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Node struct {
	ID    string
	Addr  string
	Store *KeyValueStore
	Peers []string
	Ring  *ConsistentHashRing
}

type KVPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (n *Node) forwardRequest(targetNode string, method string, url string, body io.Reader) (*http.Response, error) {
	// Create an HTTP client to forward the request to the target node
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("http://%s%s", targetNode, url), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
func (n *Node) HandlePut(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var payload KVPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil || payload.Key == "" || payload.Value == "" {
		http.Error(w, fmt.Sprintf("Invalid payload: %v", err), http.StatusBadRequest)
		return
	}

	// Find the target node using consistent hashing
	targetNode := n.Ring.GetNode(payload.Key)

	// If this node is the target, store it locally
	if targetNode == n.ID {
		n.Store.Put(payload.Key, payload.Value)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}

	// Otherwise, forward the request to the target node
	resp, err := n.forwardRequest(targetNode, "PUT", "/put", bytes.NewReader(bodyBytes))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to forward request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	w.Write(respBody)
}

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

	// Find the target node using consistent hashing
	targetNode := n.Ring.GetNode(key)

	// If this node is the target, get the value locally
	if targetNode == n.ID {
		value, ok := n.Store.Get(key)
		if !ok {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		w.Write([]byte(value)) // Just return the value
		return
	}

	// Otherwise, forward the request to the target node
	response, err := n.forwardRequest(targetNode, "GET", "/get?key="+key, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to forward request: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Relay the response from the target node
	body, _ := io.ReadAll(response.Body)
	w.Write(body)
}

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

	// Find the target node using consistent hashing
	targetNode := n.Ring.GetNode(key)

	// If this node is the target, delete it locally
	if targetNode == n.ID {
		ok := n.Store.Delete(key)
		if !ok {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusOK) // Just return HTTP 200 status
		w.Write([]byte("OK")) 
		return
	}

	// Otherwise, forward the request to the target node
	response, err := n.forwardRequest(targetNode, "DELETE", "/delete?key="+key, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to forward request: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	// Relay the response from the target node
	body, _ := io.ReadAll(response.Body)
	w.Write(body)
}

// HandleHealth provides a simple health check endpoint
func (n *Node) HandleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "[%s] OK", n.ID)
}
