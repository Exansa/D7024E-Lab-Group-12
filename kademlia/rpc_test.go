package d7024e

import (
	"encoding/json"
	"net"
	"reflect"
	"testing"
)

func TestSendMessage(t *testing.T) {
	// Create a mock RPC message
	msg := new(RPC)
	msg.Type = PING
	msg.Sender = NewContact(NewRandomKademliaID(), "127.0.0.1:8000")
	msg.Receiver = NewContact(NewRandomKademliaID(), "127.0.0.1:8080")
	msg.Data.PING = "Ping!"

	// Start a mock server to receive the message
	server, err := net.ListenPacket("udp", "127.0.0.1:8080")
	if err != nil {
		t.Fatalf("Failed to start server: %v", err)
	}
	defer server.Close()

	// Send the message
	sendMessage(msg)

	// Receive the message on the server
	buf := make([]byte, 1024)
	n, _, err := server.ReadFrom(buf)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Unmarshal the message and check its contents
	var receivedMsg RPC
	err = json.Unmarshal(buf[:n], &receivedMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	if receivedMsg.Type != msg.Type {
		t.Errorf("Received message has wrong method: expected %q, got %q", msg.Type, receivedMsg.Type)
	}
	if !reflect.DeepEqual(receivedMsg.Data, msg.Data) {
		t.Errorf("Received message has wrong params: expected %v, got %v", msg.Data, receivedMsg.Data)
	}
}

/*
func TestPing(t *testing.T) {
	kademlia := NewKademlia("localhost:8000", true)
	kademlia.initNode()
	network := NewNetwork(kademlia)
	kademlia2 := NewKademlia("localhost:9998", false)
	kademlia2.initNode()
	network.Listen("localhost:9998")
	contact := NewContact(NewRandomKademliaID(), "localhost:9998")
	network.SendPingMessage(&contact)
}*/
