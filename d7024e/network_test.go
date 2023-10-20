package d7024e

import (
	"fmt"
	"testing"
	"time"
)

func TestNetwork(t *testing.T) {
	fmt.Println("TestNetwork")
}
func TestNetworkPingProcess(t *testing.T) {
	fmt.Println("TestNetworkPingProcess")
	// Create a network
	receiver := NewKademlia("127.0.0.1:6000")
	receiver.setNodeID(NewRandomKademliaID()) // Simple way to init node
	receiverContact := NewContact(receiver.ID, receiver.ADDRESS)

	sender := NewKademlia("127.0.0.1:6001")
	sender.setNodeID(NewRandomKademliaID())

	// Start listening after messages
	go receiver.Network.Listen()
	go sender.Network.Listen()
	time.Sleep(100 * time.Millisecond)

	// Send a message to the network
	sender.Network.SendPingMessage(&receiverContact)

	// Wait for the message to be handled
	time.Sleep(100 * time.Millisecond)

	res := <-sender.Network.msgBuffer

	// Check that the message was handled
	if res.Type != PONG {
		t.Fatal("Message was not properly received")
	} else {
		fmt.Println("Message was received")
	}
}

func TestNetworkPingFunction(t *testing.T) {
	fmt.Print("TestNetworkPingFunction")
	// Create a network
	receiver := NewKademlia("127.0.0.1:6002")
	receiver.setNodeID(NewRandomKademliaID()) // Simple way to init node
	receiverContact := NewContact(receiver.ID, receiver.ADDRESS)

	sender := NewKademlia("127.0.0.1:6003")
	sender.setNodeID(NewRandomKademliaID())

	// Start listening after messages
	go receiver.Network.Listen()
	go sender.Network.Listen()

	// Send a message to the network
	sender.Network.ping(5, &receiverContact) // <- Fails or gets stuck if the message is not handled
}

func TestLookupBuffer(t *testing.T) {
	fmt.Println("TestLookupBuffer")
	// Create a buffer and fill it with 10 rpcs
	buffer := LookupBuffer{}
	for i := 0; i < 10; i++ {
		rpc := RPC{}
		rpc.Sender = NewContact(NewRandomKademliaID(), "")
		rpc.Type = FIND_NODE

		buffer.Append(rpc)
	}
	rpc := RPC{}
	rpc.Sender = NewContact(NewRandomKademliaID(), "")
	rpc.Type = FIND_NODE

	buffer.Append(rpc)
	buffer.Has(rpc)
	buffer.Remove(rpc)
	buffer.Len()

}
func TestNetworkFindNode(t *testing.T) {
	receiver := NewKademlia("127.0.0.1:4000")
	receiver.setNodeID(NewRandomKademliaID()) // Simple way to init node
	receiverContact := NewContact(receiver.ID, receiver.ADDRESS)

	sender := NewKademlia("127.0.0.1:4001")
	sender.setNodeID(NewRandomKademliaID())
	senderContact := NewContact(sender.ID, sender.ADDRESS)

	go receiver.Network.Listen()
	go sender.Network.Listen()
	time.Sleep(100 * time.Millisecond)

	// Send a message to the network
	sender.Network.findNode(receiverContact.ID, &senderContact)

}
