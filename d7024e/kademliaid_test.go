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

func TestIdComparison(t *testing.T) {
	fmt.Println("TestLess")
	kID1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2 := NewKademliaID("FFFFFFFF00000000000000000000000000000001")
	kID1.Less(kID2)
	kID1.Equals(kID2)
	kID1.String()
}
