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
	// Create a network
	receiver := NewKademlia("127.0.0.1:8000")
	receiver.setNodeID(NewRandomKademliaID()) // Simple way to init node
	receiverContact := NewContact(receiver.ID, receiver.ADDRESS)

	sender := NewKademlia("127.0.0.1:8001")
	sender.setNodeID(NewRandomKademliaID())

	// Start listening after messages
	go receiver.Network.Listen()
	go sender.Network.Listen()

	// Send a message to the network
	sender.Network.SendPingMessage(&receiverContact)

	// Wait for the message to be handled
	time.Sleep(100 * time.Millisecond)

	res := <-sender.Network.msgChan

	// Check that the message was handled
	if res.Type != PONG {
		t.Fatal("Message was not properly received")
	} else {
		fmt.Println("Message was received")
	}
}

func TestNetworkPingFunction(t *testing.T) {
	// Create a network
	receiver := NewKademlia("127.0.0.1:8002")
	receiver.setNodeID(NewRandomKademliaID()) // Simple way to init node
	receiverContact := NewContact(receiver.ID, receiver.ADDRESS)

	sender := NewKademlia("127.0.0.1:8003")
	sender.setNodeID(NewRandomKademliaID())

	// Start listening after messages
	go receiver.Network.Listen()
	go sender.Network.Listen()

	// Send a message to the network
	sender.Network.ping(&receiverContact) // <- Fails or gets stuck if the message is not handled
}
