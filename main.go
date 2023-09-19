package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/http"
	"sync"
)

// KeyValue represents a key-value pair
type KeyValue struct {
	Key   string
	Value string
}

// Node represents a node in the distributed system
type Node struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewNode() *Node {
	return &Node{
		data: make(map[string]string),
	}
}

// DistributedKeyValueStore represents the distributed key-value store
type DistributedKeyValueStore struct {
	nodes []*Node
	mu    sync.RWMutex
}

func NewDistributedKeyValueStore() *DistributedKeyValueStore {
	return &DistributedKeyValueStore{}
}

// Put stores a key-value pair in the distributed key-value store
func (d *DistributedKeyValueStore) Put(key, value string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Implement logic to distribute the data across nodes
	// You may use consistent hashing or other partitioning techniques

	// For simplicity, we'll use a round-robin approach here
	node := d.getNodeForKey(key)
	node.mu.Lock()
	defer node.mu.Unlock()
	node.data[key] = value
}

// Get retrieves the value associated with a key from the distributed key-value store
func (d *DistributedKeyValueStore) Get(key string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

	// Implement logic to locate the correct node for key retrieval

	// For simplicity, we'll use a round-robin approach here
	node := d.getNodeForKey(key)
	node.mu.RLock()
	defer node.mu.RUnlock()
	value, ok := node.data[key]
	return value, ok
}

// Delete removes a key-value pair from the distributed key-value store
func (d *DistributedKeyValueStore) Delete(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Implement logic to locate the correct node for deletion

	// For simplicity, we'll use a round-robin approach here
	node := d.getNodeForKey(key)
	node.mu.Lock()
	defer node.mu.Unlock()
	delete(node.data, key)
}

// Replicate replicates data across nodes for fault tolerance
func (d *DistributedKeyValueStore) Replicate(newNode *Node) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// Implement logic to distribute existing key-value pairs to the new node
	for _, node := range d.nodes {
		node.mu.Lock()
		for key, value := range node.data {
			newNode.data[key] = value
		}
		node.mu.Unlock()
	}

	// Add the new node to the list of nodes
	d.nodes = append(d.nodes, newNode)
}

// getNodeForKey selects the next node in a round-robin fashion
func (d *DistributedKeyValueStore) getNodeForKey(key string) *Node {
	// Use a hash function to generate a hash value from the key
	hash := fnv.New32a()
	hash.Write([]byte(key))
	hashValue := hash.Sum32()
	// Find the node responsible for the hash value
	index := int(hashValue) % len(d.nodes)
	return d.nodes[index]
}

// HandlePutRequest handles PUT requests to store key-value pairs.
func HandlePutRequest(kv *DistributedKeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		kv.Put(request.Key, request.Value)
		w.WriteHeader(http.StatusNoContent)
	}
}

// HandleGetRequest handles GET requests to retrieve values by key.
func HandleGetRequest(kv *DistributedKeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value, exists := kv.Get(key)
		if !exists {
			http.Error(w, "Key not found", http.StatusNotFound)
			return
		}

		response := struct {
			Value string `json:"value"`
		}{Value: value}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// HandleDeleteRequest handles DELETE requests to remove key-value pairs.
func HandleDeleteRequest(kv *DistributedKeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		kv.Delete(key)
		w.WriteHeader(http.StatusNoContent)
	}
}

// HandleReplicateRequest handles replicate requests to store all data in new node.
func HandleReplicateRequest(kv *DistributedKeyValueStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		newNode := NewNode()
		kv.Replicate(newNode)
		w.WriteHeader(http.StatusNoContent)
	}
}

func main() {
	distributedStore := NewDistributedKeyValueStore()
	node1 := NewNode()
	node2 := NewNode()
	node3 := NewNode()
	distributedStore.nodes = append(distributedStore.nodes, node1, node2, node3)

	http.HandleFunc("/put", HandlePutRequest(distributedStore))
	http.HandleFunc("/get", HandleGetRequest(distributedStore))
	http.HandleFunc("/delete", HandleDeleteRequest(distributedStore))
	http.HandleFunc("/replicate", HandleReplicateRequest(distributedStore))

	if err := http.ListenAndServe(":8081", nil); err != nil {
		fmt.Println(err)
	}
}
