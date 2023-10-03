package d7024e

import (
	"fmt"
	"testing"
)

func TestKademlia(t *testing.T) {
	fmt.Println("TestKademlia")
}

func TestNewKademlia(t *testing.T) {
	kademlia, _ := NewKademlia("localhost:9999", true)
	fmt.Println(kademlia.ID.String())
}
