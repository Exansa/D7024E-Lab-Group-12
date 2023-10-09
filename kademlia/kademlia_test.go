package d7024e

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"
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

func TestInitNode(t *testing.T) {
	fmt.Println("TestNodeInit2")
	rootNode := NewKademlia("localhost:1337", true) //Bootstrap route
	rootNode.initNode()
	//wait for root node to be initialized
	time.Sleep(100 * time.Millisecond)

	childNode := NewKademlia("localhost:8000", false) //Non-bootstrap route
	childNode.initNode()
	fmt.Println(childNode.isInitialized())
}

func TestStoreValue(t *testing.T) {
	fmt.Println("TestStore")
	kademlia := NewKademlia("localhost:9999", true)
	kademlia.initNode()
	dataHash := hashData([]byte("test"))
	dataKey := hex.EncodeToString(hashData([]byte("test")))
	kademlia.StoreValue(dataHash, dataKey)
}
