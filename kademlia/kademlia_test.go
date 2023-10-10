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
	kademlia := NewKademlia("127.0.0.1:9999", true)
	if kademlia.ID != nil || kademlia.RoutingTable != nil || kademlia.Network != nil {
		t.Fail()
	}
}

func TestInitNode(t *testing.T) {
	fmt.Println("TestNodeInit2")
	rootNode := NewKademlia("127.0.0.1:1337", true) //Bootstrap route
	rootNode.initNode()
	//wait for root node to be initialized
	time.Sleep(100 * time.Millisecond)

	childNode := NewKademlia("127.0.0.1:7999", false) //Non-bootstrap route
	childNode.initNode()
	time.Sleep(100 * time.Millisecond)

	fmt.Println(childNode.isInitialized())
	fmt.Println(childNode.isBootstrapNode())

	fmt.Println(rootNode.isInitialized())
	fmt.Println(rootNode.isBootstrapNode())

	if !childNode.isInitialized() || childNode.isBootstrapNode() {
		t.Fail()
	}

	if !rootNode.isInitialized() || !rootNode.isBootstrapNode() {
		t.Fail()
	}

}

func TestStoreValue(t *testing.T) {
	fmt.Println("TestStore")
	kademlia := NewKademlia("127.0.0.1:9999", true)
	kademlia.initNode()
	dataHash := hashData([]byte("test"))
	dataKey := hex.EncodeToString(hashData([]byte("test")))
	kademlia.StoreValue(dataHash, dataKey)
}
