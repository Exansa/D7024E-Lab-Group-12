package d7024e

import (
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

// func TestInitNode(t *testing.T) {
// 	fmt.Println("TestNodeInit2")
// 	rootNode := NewKademlia("127.0.0.1:1337") //Bootstrap route
// 	rootNode.initNode()
// 	//wait for root node to be initialized
// 	time.Sleep(100 * time.Millisecond)

// 	childNode := NewKademlia("127.0.0.1:7999") //Non-bootstrap route
// 	childNode.initNode()
// 	time.Sleep(100 * time.Millisecond)

// 	fmt.Println(childNode.isInitialized())
// 	fmt.Println(childNode.isBootstrapNode())

// 	fmt.Println(rootNode.isInitialized())
// 	fmt.Println(rootNode.isBootstrapNode())

// 	if !childNode.isInitialized() || childNode.isBootstrapNode() {
// 		t.Fail()
// 	}

// 	if !rootNode.isInitialized() || !rootNode.isBootstrapNode() {
// 		t.Fail()
// 	}

// }

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
	time.Sleep(300 * time.Millisecond)

	child1 := NewKademlia("127.0.0.1:7990")
	child1.initNode()
	time.Sleep(300 * time.Millisecond)

	child2 := NewKademlia("127.0.0.1:7992")
	child2.initNode()

	time.Sleep(300 * time.Millisecond)

	child3 := NewKademlia("127.0.0.1:7993")
	child3.initNode()

	time.Sleep(300 * time.Millisecond)

	child4 := NewKademlia("127.0.0.1:7994")
	child4.initNode()

	time.Sleep(300 * time.Millisecond)

	child8 := NewKademlia("127.0.0.1:7995")
	child8.initNode()

	time.Sleep(300 * time.Millisecond)

	child5 := NewKademlia("127.0.0.1:7996")
	child5.initNode()

	time.Sleep(300 * time.Millisecond)

	child6 := NewKademlia("127.0.0.1:7997")
	child6.initNode()

	time.Sleep(300 * time.Millisecond)

	child7 := NewKademlia("127.0.0.1:7998")
	child7.initNode()

	time.Sleep(300 * time.Millisecond)

	time.Sleep(1 * time.Second)

	res := child7.Network.Kademlia.LookupContact(child1.ID, child7.RoutingTable.me)
	fmt.Println("===========================================================================")
	fmt.Println("Found some contacts!")
	if res.Has(child1.ID) {
		fmt.Println("Found target contact")
		return
	}

	fmt.Println("Did not find target contact")
	fmt.Printf("Found: %s\n", res.Contacts[0].ID)
	fmt.Printf("Target: %s\n", child1.ID)
	fmt.Printf("Full list:%v\n", res.Contacts)

	fmt.Println("Root node contacts:")
	contacts := rootNode.RoutingTable.FindClosestContacts(child1.ID, 3)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
	t.Fail()

}

func TestStoreValue(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("TestStore")
	rootNode := NewKademlia("127.0.0.1:1337")
	rootNode.initNode()
	//wait for root node to be initialized
	time.Sleep(500 * time.Millisecond)

	childNode := NewKademlia("127.0.0.1:7989")
	childNode.initNode()
	time.Sleep(500 * time.Millisecond)

	child1 := NewKademlia("127.0.0.1:7990")
	child1.initNode()
	time.Sleep(500 * time.Millisecond)

	child2 := NewKademlia("127.0.0.1:7992")
	child2.initNode()

	time.Sleep(500 * time.Millisecond)

	child3 := NewKademlia("127.0.0.1:7993")
	child3.initNode()

	time.Sleep(500 * time.Millisecond)

	child4 := NewKademlia("127.0.0.1:7994")
	child4.initNode()

	time.Sleep(500 * time.Millisecond)

	child8 := NewKademlia("127.0.0.1:7995")
	child8.initNode()

	time.Sleep(500 * time.Millisecond)

	child5 := NewKademlia("127.0.0.1:7996")
	child5.initNode()

	time.Sleep(1000 * time.Millisecond)

	err := child3.Store([]byte("test"))

	time.Sleep(1000 * time.Millisecond)

	if err != nil {
		t.Fail()
	}
}

func TestFindValue(t *testing.T) {
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("TestStore")
	rootNode := NewKademlia("127.0.0.1:1337")
	rootNode.initNode()
	//wait for root node to be initialized
	time.Sleep(500 * time.Millisecond)

	childNode := NewKademlia("127.0.0.1:7989")
	childNode.initNode()
	time.Sleep(500 * time.Millisecond)

	child1 := NewKademlia("127.0.0.1:7990")
	child1.initNode()
	time.Sleep(500 * time.Millisecond)

	child2 := NewKademlia("127.0.0.1:7992")
	child2.initNode()

	time.Sleep(500 * time.Millisecond)

	child3 := NewKademlia("127.0.0.1:7993")
	child3.initNode()

	time.Sleep(500 * time.Millisecond)

	child4 := NewKademlia("127.0.0.1:7994")
	child4.initNode()

	time.Sleep(500 * time.Millisecond)

	child8 := NewKademlia("127.0.0.1:7995")
	child8.initNode()

	time.Sleep(500 * time.Millisecond)

	child5 := NewKademlia("127.0.0.1:7996")
	child5.initNode()

	time.Sleep(1000 * time.Millisecond)

	err := child3.Store([]byte("test"))

	time.Sleep(1000 * time.Millisecond)

	if err != nil {
		t.Fail()
	}
	result := child5.GetData("a94a8fe5ccb19ba61c4c0873d391e987982fbbd3")
	fmt.Println("Result ", string(result))
	res := child5.LookupData([]byte("test"))
	if string(res) == string([]byte("test")) {
		fmt.Println("Success")
	} else {
		t.Fail()
	}

}

func TestGetLocalData(t *testing.T) {
	kademlia := NewKademlia("127.0.0.1:1337")
	kademlia.initNode()
	kademlia.StoreLocally([]byte("test"), string(hashData([]byte("test"))))
	testValue := kademlia.GetData(string(hashData([]byte("test"))))
	fmt.Print("Test value: ", string(testValue), "\n")
}
