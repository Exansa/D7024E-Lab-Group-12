package d7024e

import (
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

func TestInitChildNode(t *testing.T) {
	fmt.Println("TestNodeInit2")
	kademlia1 := NewKademlia("localhost:9998", true)
	kademlia1.initNode()
	kademlia2 := NewKademlia("localhost:9999", false)
	kademlia2.initNode()
	fmt.Println(kademlia2.isInitialized())
}
