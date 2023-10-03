package d7024e

import (
	"fmt"
	"testing"
)

func TestKademliaID(t *testing.T) {
	fmt.Println("TestKademliaID")
}

func TestNewKademlia(t *testing.T) {
	kademlia := NewKademlia("10.0.0.2:9999", true)
	fmt.Println(kademlia.ADDRESS)
}
