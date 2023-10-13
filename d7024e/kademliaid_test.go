package d7024e

import (
	"fmt"
	"testing"
)

func TestKademliaID(t *testing.T) {
	fmt.Println("TestKademliaID")
}

func TestNewKademliaID(t *testing.T) {
	fmt.Println("TestNewKademliaID")
	kID := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	fmt.Println(kID.String())
}

func TestNewRandomKademliaID(t *testing.T) {
	fmt.Println("TestNewRandomKademliaID")
	kID := NewRandomKademliaID()
	fmt.Println(kID.String())
}
