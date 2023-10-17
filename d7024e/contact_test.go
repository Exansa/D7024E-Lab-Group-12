package main

import (
	"fmt"
	"testing"
)

func TestNewContact(t *testing.T) {
	newID := NewRandomKademliaID()
	newContact := NewContact(newID, "10.0.0.2")
	fmt.Println(newContact)
}

func TestCalcDistance(t *testing.T) {
	newID1 := NewRandomKademliaID()
	newID2 := NewRandomKademliaID()
	newContact1 := NewContact(newID1, "10.0.0.2")
	newContact2 := NewContact(newID2, "10.0.0.3")
	newContact1.CalcDistance(newContact2.ID)
	if newContact1.Distance == nil {
		t.Errorf("Distance is nil")
	} else {
		fmt.Println("Distance is not nil")
	}
}

func TestLess(t *testing.T) {
	newID1 := NewRandomKademliaID()
	newID2 := NewRandomKademliaID()
	newID3 := NewRandomKademliaID()
	newContact1 := NewContact(newID1, "10.0.0.2")
	newContact2 := NewContact(newID2, "10.0.0.3")
	newContact3 := NewContact(newID3, "10.0.0.4")
	newContact1.CalcDistance(newContact2.ID)
	newContact2.CalcDistance(newContact1.ID)
	newContact3.CalcDistance(newContact1.ID)
	if newContact3.Less(&newContact2) {
		fmt.Println("Contact1 is less than contact2")
	} else if newContact2.Less(&newContact3) {
		fmt.Println("Contact2 is less than contact1")
	} else {
		t.Errorf("Contact1 is not less than contact2")
	}
}

func TestString(t *testing.T) {
	newID := NewRandomKademliaID()
	newContact := NewContact(newID, "10.0.0.2")
	testString := newContact.String()
	if testString == "" {
		t.Errorf("String is empty")
	} else {
		fmt.Println("String is not empty", testString)
	}
}

func TestAppend(t *testing.T) {
	contactCandidates := new(ContactCandidates)
	newID := NewRandomKademliaID()
	newContact := NewContact(newID, "10.0.0.2")
	contactCandidates.Append([]Contact{newContact})
	if contactCandidates.Len() != 1 {
		t.Errorf("ContactCandidates length is not 1")
	} else {
		fmt.Println("ContactCandidates length is 1")
	}
}

func TestGetContacts(t *testing.T) {
	contactCandidates := new(ContactCandidates)
	newID := NewRandomKademliaID()
	newContact := NewContact(newID, "10.0.0.2")
	contactCandidates.Append([]Contact{newContact})
	result := contactCandidates.GetContacts(1)
	if len(result) != 1 {
		t.Errorf("Couldn't get contact number")
	} else {
		fmt.Println("Could get contact number")
	}
}

func TestContactLen(t *testing.T) {
	contactCandidates := new(ContactCandidates)
	newID1 := NewRandomKademliaID()
	newID2 := NewRandomKademliaID()
	newID3 := NewRandomKademliaID()
	newContact1 := NewContact(newID1, "10.0.0.2")
	newContact2 := NewContact(newID2, "10.0.0.3")
	newContact3 := NewContact(newID3, "10.0.0.4")
	contactCandidates.Append([]Contact{newContact1, newContact2, newContact3})
	if contactCandidates.Len() != 3 {
		fmt.Println("ContactCandidates length is not 3")
	} else {
		fmt.Println("ContactCandidates length is 3")
	}
}
