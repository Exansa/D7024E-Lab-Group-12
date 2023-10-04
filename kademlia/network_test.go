package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestNetwork(t *testing.T) {
	fmt.Println("TestNetwork")
}
func TestNetworkListen(t *testing.T) {
	// Create a network
	kademlia := NewKademlia("localhost:8000", true)
	network := NewNetwork(kademlia)

	// Start listening on a random port
	go network.Listen("localhost:0")

	// Send a message to the network
	msg := RPC{}
	msgBytes, err := json.Marshal(msg)

	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	conn, err := net.Dial("udp", "localhost:8000")

	if err != nil {
		t.Fatalf("Failed to dial network: %v", err)
	}
	_, err = conn.Write(msgBytes)

	if err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Wait for the message to be handled
	time.Sleep(100 * time.Millisecond)

	// Check that the message was handled
	if network.msgChan == nil {
		t.Fatal("Message was not received")
	} else {
		fmt.Println("Message was received")
	}
}

/*
func TestNetworkListen2(t *testing.T) {
	conn1, conn2 := net.Pipe()
	defer conn1.Close()
	defer conn2.Close()

	// Create a networks
	kademlia := NewKademlia(conn1.LocalAddr().String()+":0", true)
	network := NewNetwork(kademlia)

	// Create a pipe

	// Start listening on the pipe
	go network.Listen(conn2.LocalAddr().String() + ":0")

	// Send a message through the pipe
	msg := RPC{}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		t.Fatalf("Failed to marshal message: %v", err)
	}
	_, err = conn1.Write(msgBytes)
	if err != nil {
		t.Fatalf("Failed to write message: %v", err)
	}

	// Read the message from the other end of the pipe
	buf := make([]byte, 1024)
	n, err := conn2.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	// Unmarshal the message and check that it's correct
	var receivedMsg RPC
	err = json.Unmarshal(buf[:n], &receivedMsg)
	if err != nil {
		t.Fatalf("Failed to unmarshal message: %v", err)
	}
	if !reflect.DeepEqual(msg, receivedMsg) {
		t.Fatal("Received message does not match sent message")
	}
}
*/
