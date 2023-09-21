package main

import (
	"bufio"
	"fmt"
	"hash/fnv"
	"log"
	"os"
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
func (d *DistributedKeyValueStore) Put(key, value string) (string, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()

	node := d.getNodeForKey(key)
	node.mu.Lock()
	defer node.mu.Unlock()
	node.data[key] = value
	value, ok := node.data[key]
	return value, ok

}

// Get retrieves the value associated with a key from the distributed key-value store
func (d *DistributedKeyValueStore) Get(key string) (string, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()

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

func main() {
	distributedStore := NewDistributedKeyValueStore()
	node1 := NewNode()
	node2 := NewNode()
	node3 := NewNode()
	distributedStore.nodes = append(distributedStore.nodes, node1, node2, node3)
	for {
		fmt.Println("Select an option:")
		fmt.Println("1. Put")
		fmt.Println("2. Get")
		fmt.Println("3. Delete")
		fmt.Println("4. Replicate")
		fmt.Println("5. Quit")

		var choice int

		_, err := fmt.Scan(&choice)

		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 1 and 5.")
			continue
		}

		switch choice {
		case 1:
			fmt.Println("You selected Option Put. Running code for Option 1...")
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter key: ")
			key, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading data: %s\n", err)
			}
			fmt.Print("Enter Value: ")

			value, errValue := reader.ReadString('\n')
			if errValue != nil {
				log.Fatalf("Error reading data: %s\n", errValue)
			}
			_, ok := distributedStore.Put(key, value)
			if !ok {
				log.Fatalf("Error while Put")
			}
			break

		case 2:
			fmt.Println("You selected Option Get. Running code for Option 2...")
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter key: ")
			key, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading data: %s\n", err)
			}
			value, exists := distributedStore.Get(key)
			if !exists {
				log.Fatalf("Key not found")
				return
			}

			response := struct {
				Key   string `json:"key"`
				Value string `json:"value"`
			}{Key: key, Value: value}

			fmt.Printf("%s:%s", response.Key, response.Value)
			break

		case 3:
			fmt.Println("You selected Option Delete. Running code for Option 3...")
			reader := bufio.NewReader(os.Stdin)

			fmt.Print("Enter key: ")
			key, err := reader.ReadString('\n')
			if err != nil {
				log.Fatalf("Error reading data: %s\n", err)
			}
			distributedStore.Delete(key)
			break

		case 4:
			fmt.Println("You selected Option Replicate. Running code for Option 4...")
			newNode := NewNode()
			distributedStore.Replicate(newNode)
			fmt.Println("New Node data:", newNode.data)
			break

		case 5:
			fmt.Println("Exiting the program.")
			os.Exit(0)
			break
		default:
			fmt.Println("Invalid choice. Please select a valid option.")
		}
	}
}
