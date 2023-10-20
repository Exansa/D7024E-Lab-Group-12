package d7024e

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8000")
	rt := NewRoutingTable(&contact)

	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "127.0.0.1:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "127.0.0.1:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "127.0.0.1:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "127.0.0.1:8002"))
	rt.AddContact(NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "127.0.0.1:8002"))

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
	rt.getBucketIndex(NewKademliaID("3111111400000000000000000000000000000000"))
}
