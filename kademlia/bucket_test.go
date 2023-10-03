package d7024e

import (
	"fmt"
	"testing"
)

func TestNewBucket(t *testing.T) {
	bucket := newBucket()
	fmt.Println(bucket)
}

func TestLen(t *testing.T) {
	bucket := newBucket()
	fmt.Println(bucket.Len())
}

func TestAddContact(t *testing.T) {
	bucket := newBucket()
	newID1 := NewRandomKademliaID()
	contactInfo1 := NewContact(newID1, "10.0.0.2")
	bucket.AddContact(contactInfo1)
	newID2 := NewRandomKademliaID()
	contactInfo2 := NewContact(newID2, "10.0.0.3")
	bucket.AddContact(contactInfo2)
	if bucket.Len() != 2 {
		t.Errorf("Bucket length is not 2")
	} else {
		fmt.Println("Bucket length is 2")
	}
}
