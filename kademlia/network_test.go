package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"reflect"
	"testing"
)

func TestNetwork(t *testing.T) {
	fmt.Println("TestNetwork")
}

func TestNetworkListen(t *testing.T) {
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
