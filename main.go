package main

import (
	"distributed-key-value-store/node"
	"log"
	"net/http"
	"sync"
)

// loggingMiddleware logs incoming requests
func loggingMiddleware(n *node.Node, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s %s", n.ID, r.Method, r.URL.Path)
		next(w, r)
	}
}

func StartServer(addr string, store *node.KeyValueStore, wg *sync.WaitGroup) {
	defer wg.Done() // Signal when the server has started

	// List of all nodes for the consistent hash ring
	allNodes := []string{"localhost:8080", "localhost:8081", "localhost:8082"}

	// Initialize the consistent hash ring with replicas
	ring := node.NewConsistentHashRing(3)
	ring.AddNodes(allNodes)

	// Create a new node
	n := &node.Node{
		ID:    "localhost" + addr, // ID of the node
		Addr:  addr,
		Store: store,    // The store passed in for this node
		Peers: allNodes, // List of all peers
		Ring:  ring,     // Consistent hash ring
	}

	// Create a new ServeMux for each server
	mux := http.NewServeMux()

	// Define the routes for GET, PUT, DELETE, and HEALTH operations
	mux.HandleFunc("/get", loggingMiddleware(n, n.HandleGet))
	mux.HandleFunc("/put", loggingMiddleware(n, n.HandlePut))
	mux.HandleFunc("/delete", loggingMiddleware(n, n.HandleDelete))
	mux.HandleFunc("/health", loggingMiddleware(n, n.HandleHealth))

	// Start the server and listen on the specified address
	log.Printf("Starting node at %s", n.ID)
	if err := http.ListenAndServe(n.Addr, mux); err != nil {
		log.Fatalf("Error starting server at %s: %v", addr, err)
	}
}

func main() {
	var wg sync.WaitGroup

	// Start three nodes on different ports
	wg.Add(3) // We have 3 nodes to start
	go StartServer(":8080", node.NewKeyValueStore(), &wg)
	go StartServer(":8081", node.NewKeyValueStore(), &wg)
	go StartServer(":8082", node.NewKeyValueStore(), &wg)

	// Wait for all goroutines (servers) to finish
	wg.Wait()
}
