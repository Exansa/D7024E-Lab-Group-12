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

func TestUpdateIDParams(t *testing.T) {
	fmt.Println("TestUpdateIDParams")
	kademlia := NewKademlia("127.0.0.1:7888")
	newID := NewKademliaID(bootstrapIDString)
	kademlia.updateIDParams(newID)
	if !kademlia.ID.Equals(newID) {
		t.Fail()
	}
}

func TestLookupContact(t *testing.T) {
	fmt.Println("TestLookupContact")
	rootNode := NewKademlia("127.0.0.1:1337")
	rootNode.initNode()
	//wait for root node to be initialized
	time.Sleep(100 * time.Millisecond)

	child1 := NewKademlia("127.0.0.1:7990")
	child1.initNode()

	child2 := NewKademlia("127.0.0.1:7992")
	child2.initNode()

	child3 := NewKademlia("127.0.0.1:7993")
	child3.initNode()

	child4 := NewKademlia("127.0.0.1:7994")
	child4.initNode()

	child8 := NewKademlia("127.0.0.1:7995")
	child8.initNode()

	child5 := NewKademlia("127.0.0.1:7996")
	child5.initNode()

	child6 := NewKademlia("127.0.0.1:7997")
	child6.initNode()

	child7 := NewKademlia("127.0.0.1:7998")
	child7.initNode()

	time.Sleep(100 * time.Millisecond)

	res := child6.Network.Kademlia.LookupContact(child1.ID)
	if !res.contacts[0].ID.Equals(child3.ID) {
		fmt.Println("Did not find contact")
		fmt.Printf("Found: %s\n", res.contacts[0].ID)
		fmt.Printf("Target: %s\n", child3.ID)
		fmt.Printf("Full list:%v\n", res.contacts)
		fmt.Printf("Root node contacts: %v\n", rootNode.RoutingTable.buckets[1].Len())

		contacts := rootNode.RoutingTable.FindClosestContacts(child1.ID, 3)
		for i := range contacts {
			fmt.Println(contacts[i].String())
		}
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
