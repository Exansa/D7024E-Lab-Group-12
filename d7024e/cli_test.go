package d7024e

import (
	"bytes"
	"testing"
	"time"
)

var buf bytes.Buffer

// func TestCLIPing(t *testing.T) {
// 	fmt.Println("TestCLI")
// 	node := NewKademlia("10.0.0.2")

// 	buf.Write([]byte("ping 10.0.0.2"))
// 	CLI(&buf, node)
// }

// func TestCLIget(t *testing.T) {
// 	fmt.Println("TestCLI")
// 	node := NewKademlia("10.0.0.2")
// 	buf.Write([]byte("get"))
// 	CLI(&buf, node)
// }

// func TestCLIPut(t *testing.T) {
// 	fmt.Println("TestCLI")
// 	node := NewKademlia("10.0.0.2")
// 	buf.Write([]byte("put image.jpg"))
// 	CLI(&buf, node)
// }

// func TestCLIExit(t *testing.T) {
// 	fmt.Println("TestCLI")
// 	node := NewKademlia("10.0.0.2")
// 	buf.Write([]byte("exit"))
// 	CLI(&buf, node)
// }

// func TestCLIHelp(t *testing.T) {
// 	fmt.Println("TestCLI")
// 	node := NewKademlia("10.0.0.2")
// 	buf.Write([]byte("help"))
// 	CLI(&buf, node)
// }

func TestCLI(t *testing.T) {
	// Create a new Kademlia instance
	kademlia := NewKademlia("127.0.0.1:5000")
	kademlia.setNodeID(NewRandomKademliaID())
	kademlia2 := NewKademlia("127.0.0.1:5001")
	kademlia2.setNodeID(NewRandomKademliaID())

	go kademlia.Network.Listen()
	go kademlia2.Network.Listen()
	time.Sleep(100 * time.Millisecond)
	// Create a new buffer for stdin

	buf.Write([]byte("help"))
	go CLI(&buf, kademlia)
	inputPing := []string{"ping", "127.0.0.1:5001"}
	ping(inputPing, kademlia)
	inputExit := []string{"test"}
	exit(inputExit, kademlia)
}
