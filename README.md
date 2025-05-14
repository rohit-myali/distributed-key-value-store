# Distributed Key-Value Store with Consistent Hashing

## Overview

This project implements a distributed key-value store in Go, utilizing **consistent hashing** to distribute keys across multiple nodes. Each node in the system stores a subset of the data and handles requests concurrently. The system is designed to be scalable and fault-tolerant, with health check endpoints for monitoring the nodes. Each node exposes CRUD (Create, Read, Update, Delete) operations, and the system ensures that data consistency is maintained.

---

## Features

* **Consistent Hashing**: Data is distributed across nodes based on a consistent hashing ring, minimizing the reorganization of data when nodes are added or removed.
* **Concurrency**: Supports concurrent handling of read (`GET`), write (`PUT`), and delete (`DELETE`) operations across multiple nodes.
* **Fault Tolerance**: Health check endpoints are provided to monitor node status and ensure the system's resilience.
* **Simple CRUD API**: Exposes a basic API for interacting with the key-value store, with endpoints for `PUT`, `GET`, `DELETE`, and `HEALTH` checks.

---

## Design Decisions

### 1. **Key-Value Store**

* Each node in the system has its own in-memory key-value store that handles CRUD operations.
* The store supports the basic operations of `PUT`, `GET`, and `DELETE`.

### 2. **Consistent Hashing**

* Each node is assigned a part of the key-space using consistent hashing.
* This allows efficient distribution of keys across the nodes and ensures that minimal data is reshuffled when nodes are added or removed.
* The system uses a hash ring to map keys to nodes.

### 3. **Concurrency**

* The system uses Go's goroutines to handle incoming HTTP requests concurrently. Each node handles requests independently, enabling high throughput and responsiveness.

### 4. **Fault Tolerance**

* The system includes basic health checks (/health) to monitor node availability. However, full fault tolerance (e.g., automatic failure handling, data replication, and recovery) is not yet implemented.
*TODO*: Add support for automatic node failure detection, data replication, and graceful recovery.


### 5. **Node Setup**

* Multiple nodes are created and run on different ports, with each node maintaining a unique address and key-value store.

---

## System Architecture

### 1. **Node**:

* Each node in the system has a unique ID and address.
* A node consists of a key-value store, a consistent hash ring, and a list of peer nodes in the cluster.

### 2. **Key-Value Store**:

* A simple in-memory key-value store is used to handle CRUD operations.
* Each node has its own key-value store, which it accesses for operations on keys it is responsible for.

### 3. **Consistent Hashing**:

* A consistent hashing ring is created to map keys to nodes based on the hash of the key.
* Each key is directed to the correct node based on the hash ring.

### 4. **API Endpoints**:

* **`/get`**: Retrieves the value of a given key from the appropriate node.
* **`/put`**: Stores a key-value pair in the appropriate node.
* **`/delete`**: Deletes a key from the appropriate node.
* **`/health`**: Provides a health check for the node.

### 5. **Health Checking**:

* Each node exposes a simple health check endpoint to ensure the node is operational.
* The health check returns an "OK" message if the node is functioning correctly.

---

## Deployment Instructions

### Prerequisites

* Go 1.18+ installed.
* Basic understanding of Go and HTTP servers.
* Multiple machines or ports to run multiple nodes (or use Docker to simulate nodes on a single machine).

### Running the System

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/your-repository/distributed-key-value-store.git
   cd distributed-key-value-store
   ```

2. **Build the Project**:

   ```bash
   go build -o distributed-key-value-store
   ```

3. **Run the Nodes**:

   The `main` function in this project starts three nodes, each listening on a different port (`8080`, `8081`, `8082`). To run the nodes, simply execute the built executable.

   ```bash
   ./distributed-key-value-store
   ```

   This will start three nodes on the following ports:

   * Node 1: `localhost:8080`
   * Node 2: `localhost:8081`
   * Node 3: `localhost:8082`

4. **Test the System**:

   You can perform `PUT`, `GET`, and `DELETE` operations using `curl`:

   * **Insert Data**:

     ```bash
     curl -X PUT http://localhost:8080/put -d '{"key":"alpha", "value":"value1"}'
     curl -X PUT http://localhost:8080/put -d '{"key":"bravo", "value":"value2"}'
     curl -X PUT http://localhost:8080/put -d '{"key":"charlie", "value":"value3"}'
     ```

   * **Get Data**:

     ```bash
     curl http://localhost:8080/get?key=alpha
     curl http://localhost:8080/get?key=bravo
     curl http://localhost:8080/get?key=charlie
     ```

   * **Delete Data**:

     ```bash
     curl -X DELETE http://localhost:8080/delete?key=alpha
     curl -X DELETE http://localhost:8080/delete?key=bravo
     curl -X DELETE http://localhost:8080/delete?key=charlie
     ```

5. **Health Check**:

   You can check the health of each node by accessing the `/health` endpoint:

   ```bash
   curl http://localhost:8080/health
   curl http://localhost:8081/health
   curl http://localhost:8082/health
   ```

   If the node is running correctly, it will return an "OK" message.

---

## Future Enhancements

* **Replication**: Implement data replication to ensure high availability and fault tolerance.
* **Node Failover**: Implement automatic failover mechanisms to recover from node failures.
* **Persistent Storage**: Switch to persistent storage (e.g., disk-based or database-backed) for data durability.
* **Load Balancing**: Implement a load balancing mechanism to handle varying node loads and traffic distribution efficiently.

---

## Conclusion

This distributed key-value store project implements a fault-tolerant and scalable system using Go, consistent hashing, and simple in-memory data stores. With concurrency support and health checks, it provides a reliable and efficient way to manage distributed key-value data across multiple nodes.

