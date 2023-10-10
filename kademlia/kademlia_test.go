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
	kademlia := NewKademlia("127.0.0.1:9999")
	if kademlia.ID != nil || kademlia.RoutingTable != nil || kademlia.Network != nil {
		t.Fail()
	}
}

func TestInitNode(t *testing.T) {
	fmt.Println("TestNodeInit2")
	rootNode := NewKademlia("127.0.0.1:1337") //Bootstrap route
	rootNode.initNode()
	//wait for root node to be initialized
	time.Sleep(100 * time.Millisecond)

	childNode := NewKademlia("127.0.0.1:7999") //Non-bootstrap route
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
	kademlia := NewKademlia("127.0.0.1:9999")
	kademlia.initNode()
	dataHash := hashData([]byte("test"))
	dataKey := hex.EncodeToString(hashData([]byte("test")))
	kademlia.StoreValue(dataHash, dataKey)
}

func TestLookupContact(t *testing.T) {
	fmt.Println("TestLookupContact")
	kademlia1 := NewKademlia("127.0.0.1:1337")
	kademlia1.initNode()
	time.Sleep(100 * time.Millisecond)
	kademlia2 := NewKademlia("127.0.0.1:1002")
	kademlia2.initNode()
	time.Sleep(100 * time.Millisecond)
	kademlia3 := NewKademlia("127.0.0.1:1003")
	kademlia3.initNode()
	time.Sleep(100 * time.Millisecond)
	kademlia4 := NewKademlia("127.0.0.1:1004")
	kademlia4.initNode()
	time.Sleep(100 * time.Millisecond)
	kademlia5 := NewKademlia("127.0.0.1:1005")
	kademlia5.initNode()
	time.Sleep(100 * time.Millisecond)
	kademlia1.RoutingTable.AddContact(*kademlia2.RoutingTable.me)
	kademlia1.RoutingTable.AddContact(*kademlia3.RoutingTable.me)
	kademlia1.RoutingTable.AddContact(*kademlia4.RoutingTable.me)
	kademlia1.RoutingTable.AddContact(*kademlia5.RoutingTable.me)
	fmt.Print("added contacts")

	shortlist := kademlia1.LookupContact(kademlia4.ID)
	if shortlist.contacts[0] != *kademlia4.RoutingTable.me {
		t.Fail()
	}
}
