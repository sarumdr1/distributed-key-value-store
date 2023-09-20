package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestKeyValueStore tests the Put , Get and Delete operations of the distributed key-value store.
func TestKeyValueStore(t *testing.T) {
	distributedStore := NewDistributedKeyValueStore()
	node1 := NewNode()
	node2 := NewNode()
	distributedStore.nodes = append(distributedStore.nodes, node1, node2)

	// Create an HTTP test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/put":
			HandlePutRequest(distributedStore)(w, r)
		case "/get":
			HandleGetRequest(distributedStore)(w, r)
		case "/delete":
			HandleDeleteRequest(distributedStore)(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	// Test PUT operation
	putData := struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{
		Key:   "test_key",
		Value: "test_value",
	}
	putDataJSON, _ := json.Marshal(putData)
	putResponse, err := http.Post(ts.URL+"/put", "application/json", bytes.NewBuffer(putDataJSON))
	if err != nil {
		t.Fatalf("PUT request failed: %v", err)
	}
	if putResponse.StatusCode != http.StatusNoContent {
		t.Errorf("PUT request returned unexpected status code: %d", putResponse.StatusCode)
	}

	// Test GET operation
	getResponse, err := http.Get(ts.URL + "/get?key=test_key")
	if err != nil {
		t.Fatalf("GET request failed: %v", err)
	}
	if getResponse.StatusCode != http.StatusOK {
		t.Errorf("GET request returned unexpected status code: %d", getResponse.StatusCode)
	}
	var getValue struct {
		Value string `json:"value"`
	}
	if err := json.NewDecoder(getResponse.Body).Decode(&getValue); err != nil {
		t.Fatalf("Failed to decode GET response: %v", err)
	}
	if getValue.Value != "test_value" {
		t.Errorf("GET request returned unexpected value: %s", getValue.Value)
	}

	// Test DELETE operation
	deleteRequest, err := http.NewRequest("DELETE", ts.URL+"/delete?key=test_key", nil)
	if err != nil {
		t.Errorf("Failed to create DELETE request: %v", err)
		return
	}
	deleteResponse, err := http.DefaultClient.Do(deleteRequest)
	if err != nil {
		t.Errorf("DELETE request failed: %v", err)
		return
	}
	if deleteResponse.StatusCode != http.StatusNoContent {
		t.Errorf("DELETE request returned unexpected status code: %d", deleteResponse.StatusCode)
	}
	// Verify that the key was deleted
	getDeletedResponse, err := http.Get(ts.URL + "/get?key=test_key")
	if err != nil {
		t.Errorf("GET request after delete failed: %v", err)
		return
	}
	if getDeletedResponse.StatusCode != http.StatusNotFound {
		t.Errorf("GET request after delete returned unexpected status code: %d", getDeletedResponse.StatusCode)
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

	value, _ := distributedStore.Put("key1", "value1")

	if value != "value1" {
		t.Errorf("Value %s is expected but got %s", "value1", value)
	}

}
