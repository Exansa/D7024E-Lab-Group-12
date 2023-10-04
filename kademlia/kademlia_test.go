package d7024e

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestKademlia(t *testing.T) {
	fmt.Println("TestKademlia")
}

func TestNewKademlia(t *testing.T) {
	kademlia := NewKademlia("localhost:9999", true)
	if kademlia.ID != nil || kademlia.RoutingTable != nil || kademlia.Network != nil {
		t.Fail()
	}
}

func TestInitRootNode(t *testing.T) {
	fmt.Println("TestNodeInit1")
	kademlia := NewKademlia("localhost:9999", true)
	kademlia.initNode()
	fmt.Println(kademlia.isInitialized())
}

/*
func TestInitNodes(t *testing.T) {
	// Create a root node
	kademlia1 := NewKademlia("localhost:9999", true)
	kademlia1.initNode()

	// Check that the root node is initialized
	if !kademlia1.isInitialized() {
		t.Fatal("Root node is not initialized")
	}

	// Create a child node
	kademlia2 := NewKademlia("localhost:9998", false)
	kademlia2.initNode()
	//network := NewNetwork(kademlia2)

	// Check that the child node is initialized
	if !kademlia2.isInitialized() {
		t.Fatal("Child node is not initialized")
	}
}*/

func TestStoreValue(t *testing.T) {
	fmt.Println("TestStore")
	kademlia := NewKademlia("localhost:9999", true)
	kademlia.initNode()
	dataHash := hashData([]byte("test"))
	dataKey := hex.EncodeToString(hashData([]byte("test")))
	kademlia.StoreValue(dataHash, dataKey)
}
