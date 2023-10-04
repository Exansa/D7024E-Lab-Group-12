package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"sync"
)

const alpha = 3
const k = 20

type Kademlia struct {
	ID           *KademliaID       //id
	ADDRESS      string            //ip:port
	DataStore    map[string][]byte //data storage
	Bootstrap    bool              //bootstrap eller inte
	RoutingTable *RoutingTable     //routingtable
	Network      *Network          //network
}

func NewKademlia(address string, bootstrap bool) *Kademlia {
	kademlia := Kademlia{}

	kademlia.ID = nil // Will get set during init
	kademlia.ADDRESS = address
	kademlia.DataStore = make(map[string][]byte)
	kademlia.Bootstrap = bootstrap //TODO: Implement logic for bootstrap here
	kademlia.RoutingTable = nil    // Will get set during init
	kademlia.Network = nil         // Will get set during init

	return &kademlia
}

func (kademlia *Kademlia) setNodeID(id *KademliaID) {
	kademlia.ID = id
	kademlia.RoutingTable = NewRoutingTable(NewContact(kademlia.ID, kademlia.ADDRESS))
	kademlia.Network = NewNetwork(kademlia)
}

// Checks if the node is initialized
//
// PANICS if the node is in an inconsistent state
func (kademlia *Kademlia) isInitialized() bool {
	if kademlia.ID != nil && kademlia.RoutingTable != nil && kademlia.Network != nil {
		return true
	} else if kademlia.ID == nil && kademlia.RoutingTable == nil && kademlia.Network == nil {
		return false
	} else {
		panic("Kademlia is in an inconsistent state")
	}
}

func (kademlia *Kademlia) initNode() {
	bootstrapAddress := "localhost:8000"
	bootstrapID := NewKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	// Check if bootstrap node
	if kademlia.Bootstrap {
		// Set static ID to bootstrap node for easy access
		kademlia.setNodeID(bootstrapID)

	} else {
		// Set a random ID
		kademlia.setNodeID(NewRandomKademliaID())

		// Find and contact bootstrap node with static ID
		bootstrapContact := NewContact(bootstrapID, bootstrapAddress)
		//TODO: Ping bootstrap node until it responds
		for {
			err := kademlia.Network.SendPingMessage(&bootstrapContact)
			if err == nil {
				// Bootstrap node is alive, connect to it
				kademlia.Network.Listen(bootstrapAddress) //TODO: Fix the listening function
				break
			} else {
				// Error log
			}
		}
	}

	// Await content updates
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {

	shortlist := ContactCandidates{}
	shortlist.contacts = kademlia.RoutingTable.FindClosestContacts(target.ID, 3)
	closest := shortlist.contacts[0]
	probed := make(map[string]bool)
	probed[closest.ID.String()] = true
	wg := sync.WaitGroup{}

	for {
		lastClosest := closest
		queue := make(chan ContactCandidates, k)

		for _, contact := range shortlist.contacts {
			// Skip if already probed to prevent dupes
			if _, ok := probed[contact.ID.String()]; ok {
				continue
			}
			probed[contact.ID.String()] = true

			wg.Add(1)

			//async FIND_NODE RPC to the closest nodes in shortlist
			go func(contact *Contact) {

				res, err := kademlia.Network.SendFindContactMessage(target, contact)

				if err != nil {
					fmt.Println("Error listening:", err.Error())
					wg.Done()
					return // If it fails to reply, it won't be added to the shortlist
				}

				queue <- res
				wg.Done()

			}(&contact)

		}

		wg.Wait()
		// Wait for all responses
		close(queue)
		for t := range queue {
			// Add all new contacts to the shortlist
			shortlist.Append(t.contacts)
		}

		// Sort shortlist
		shortlist.Sort()

		// Pick all alpha closest nodes as the new shortlist
		if shortlist.Len() > alpha {
			shortlist.contacts = shortlist.GetContacts(alpha)
		}

		// Exit the loop if no closer nodes are found
		closest = shortlist.contacts[0]
		if closest.ID.Equals(lastClosest.ID) {
			break
		}
	}
	return closest

}

func (kademlia *Kademlia) LookupData(hash string) (data string) {
	// similar to lookupcontact
	encodedData := kademlia.GetData(hash)
	//store the value in the closest node that isn't the correct node
	return string(encodedData)
}

func (kademlia *Kademlia) GetData(hash string) (data []byte) {
	return kademlia.DataStore[hash]
}

func (kademlia *Kademlia) Store(data []byte) (self bool, closest Contact, dataHash string) {
	// get hash of data
	dataHash = hex.EncodeToString(sha1.New().Sum(data))
	dataKey := NewKademliaID(dataHash)

	// find closest nodes, Maybe use lookupcontact?
	contacts := kademlia.RoutingTable.FindClosestContacts(dataKey, 3)
	closest = kademlia.RoutingTable.me
	// send store message to closest nodes
	for _, contact := range contacts {
		if contact.ID.CalcDistance(contact.ID).Less(closest.ID.CalcDistance(contact.ID)) {
			closest = contact
		}
	}
	if closest.ID.Equals(kademlia.RoutingTable.me.ID) {
		kademlia.StoreValue(data, dataHash)
		return true, closest, dataHash
	} else {
		return false, closest, dataHash
	}
	//kademlia.Network.SendStoreMessage(data, closest)
	// store data in datastore
	// kademlia.DataStore[hash] = data
}

func (kademlia *Kademlia) StoreValue(data []byte, dataHash string) {
	kademlia.DataStore[dataHash] = data
}
