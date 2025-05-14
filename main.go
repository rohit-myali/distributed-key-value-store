package main

import (
	"distributed-key-value-store/node"
	"log"
	"net/http"
)

func loggingMiddleware(n *node.Node, next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Printf("[%s] %s %s", n.ID, r.Method, r.URL.Path)
        next(w, r)
    }
}

func main() {
	// Initialize the key-value store on a single node
	n := &node.Node{
		ID:    "localhost:8080", // Node ID
		Addr:  ":8080",          // Listening on port 8080
		Store: node.NewKeyValueStore(),
	}

	// Define the routes for GET, PUT, and DELETE operations
    http.HandleFunc("/get", loggingMiddleware(n, n.HandleGet))
    http.HandleFunc("/put", loggingMiddleware(n, n.HandlePut))
    http.HandleFunc("/delete", loggingMiddleware(n, n.HandleDelete))
    http.HandleFunc("/health", loggingMiddleware(n, n.HandleHealth))

	// Start the server and listen on the specified address
	log.Printf("Starting node at %s", n.ID)
	log.Fatal(http.ListenAndServe(n.Addr, nil))
}
