package main

import (
	"testing"
)

// TestKeyValueStore tests the Put , Get and Delete operations of the distributed key-value store.
func TestKeyValueStore(t *testing.T) {
	distributedStore := NewDistributedKeyValueStore()
	node1 := NewNode()
	distributedStore.nodes = append(distributedStore.nodes, node1)
	key := "testKey"
	inputValue := "testValue"
	value, _ := distributedStore.Put(key, inputValue)

	if value != inputValue {
		t.Errorf("Value %s is expected but got %s", inputValue, value)
	}

	// Test Delete
	distributedStore.Delete(key)

	// Verify that the key no longer exists
	_, exists := distributedStore.Get(key)
	if exists {
		t.Errorf("Expected Delete to remove the key")
	}
}

func TestGet(t *testing.T) {

	distributedStore := NewDistributedKeyValueStore()
	node1 := NewNode()
	distributedStore.nodes = append(distributedStore.nodes, node1)

	distributedStore.nodes[0].data["key1"] = "value1"
	value, _ := distributedStore.Get("key1")

	if value != "value1" {
		t.Errorf("Value %s is expected but got %s", "value1", value)
	}

}

func TestPut(t *testing.T) {
	distributedStore := NewDistributedKeyValueStore()
	node1 := NewNode()
	distributedStore.nodes = append(distributedStore.nodes, node1)

	key := "testKey"
	inputValue := "testValue"
	value, _ := distributedStore.Put(key, inputValue)

	if value != inputValue {
		t.Errorf("Value %s is expected but got %s", inputValue, value)
	}
}
