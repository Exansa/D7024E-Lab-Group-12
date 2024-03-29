package main

import (
	"encoding/hex"
	"fmt"
	"sync"
)

const alpha = 3
const k = 20
const uninitIDString = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
const bootstrapIDString = "0000000000000000000000000000000000000000"

type Kademlia struct {
	ID           *KademliaID       //id
	ADDRESS      string            //ip:port
	DataStore    map[string][]byte //data storage
	RoutingTable *RoutingTable     //routingtable
	Network      *Network          //network
}

func NewKademlia(address string) *Kademlia {
	kademlia := Kademlia{}

	kademlia.ID = NewKademliaID(uninitIDString) // Placeholder until init
	kademlia.ADDRESS = address
	kademlia.DataStore = make(map[string][]byte)
	kademlia.Network = NewNetwork(&kademlia) // Will get set during init

	me := NewContact(kademlia.ID, kademlia.ADDRESS)
	kademlia.RoutingTable = NewRoutingTable(&me)

	return &kademlia
}

func (kademlia *Kademlia) setNodeID(id *KademliaID) {
	kademlia.ID = id
}

// Checks if the node is initialized
//
// PANICS if the node is in an inconsistent state

func (kademlia *Kademlia) updateIDParams(id *KademliaID) {
	kademlia.ID = id
	me := NewContact(kademlia.ID, kademlia.ADDRESS)
	kademlia.RoutingTable = NewRoutingTable(&me)
}

func (kademlia *Kademlia) isInitialized() bool {
	uninitID := NewKademliaID(uninitIDString)
	return !kademlia.ID.Equals(uninitID)
}

func (kademlia *Kademlia) isBootstrapNode() bool {
	bootstrapID := NewKademliaID(bootstrapIDString)
	return kademlia.ID.Equals(bootstrapID)
}

func (kademlia *Kademlia) initNode() error {
	bootstrapAddress := "10.0.8.3:8000"
	bootstrapID := NewKademliaID(bootstrapIDString)

	if kademlia.ADDRESS == bootstrapAddress {
		kademlia.updateIDParams(bootstrapID)
		go kademlia.Network.Listen()
		return nil
	}

	kademlia.updateIDParams(NewRandomKademliaID())

	// Check if bootstrap node is alive
	bootstrapContact := NewContact(bootstrapID, bootstrapAddress)
	go kademlia.Network.Listen()

	// Try manually pinging the bootstrap node
	err := kademlia.Network.ping(300, &bootstrapContact)
	if err == nil {
		// Bootstrap node is alive and has added you as a contact, init connection
		fmt.Println("Bootstrap node is alive, initializing connection")
		kademlia.LookupContact(kademlia.ID, kademlia.RoutingTable.me)
		return nil
	} else {
		// Invalid/No response from bootstrap node, set bootstrap node to self
		fmt.Println("Failed to attach to bootstrap node")
		return err
	}
}

func (kademlia *Kademlia) LookupContact(target *KademliaID, sender *Contact) ContactCandidates {

	shortlist := ContactCandidates{}
	contacts := kademlia.RoutingTable.FindClosestContacts(target, 3)
	shortlist.Append(contacts)

	if shortlist.Has(target) || target.Equals(kademlia.RoutingTable.me.ID) {
		return shortlist
	}

	closest := *kademlia.RoutingTable.me
	probed := ContactCandidates{}
	probed.Append([]Contact{*kademlia.RoutingTable.me})

	wg := sync.WaitGroup{}

	for {
		lastClosest := closest
		queue := make(chan ContactCandidates, k)

		if probed.Len() < alpha {
			for _, contact := range shortlist.Contacts {
				// Skip if already probed to prevent dupes
				if probed.Has(contact.ID) || contact.ID.Equals(sender.ID) {
					continue
				}

				probed.Append([]Contact{contact})

				wg.Add(1)

				//async FIND_NODE RPC to the closest nodes in shortlist
				go func(contact *Contact) {

					res, err := kademlia.Network.findNode(target, contact)

					if err != nil {
						wg.Done()
						return // If it fails to reply, it won't be added to the shortlist
					}

					queue <- res
					wg.Done()

				}(&contact)

			}
		}

		wg.Wait()
		// Wait for all responses
		close(queue)

		for t := range queue {
			// Add all new contacts to the shortlist
			shortlist.Append(t.Contacts)
		}

		// Sort shortlist
		shortlist.Sort()

		// Pick all alpha closest nodes as the new shortlist
		if shortlist.Len() > alpha {
			shortlist.Contacts = shortlist.GetContacts(alpha)
		}

		// Exit the loop if no closer nodes are found
		closest = shortlist.Contacts[0]
		if closest.ID.Equals(lastClosest.ID) || closest.ID.Equals(kademlia.RoutingTable.me.ID) || shortlist.Has(target) {
			break
		}
	}
	return shortlist

}

func (kademlia *Kademlia) LookupData(hash []byte) (data []byte) {
	dataHash := hex.EncodeToString(hashData(hash))
	dataKey := NewKademliaID(dataHash)
	shortList := kademlia.LookupContact(dataKey, kademlia.RoutingTable.me)
	// time.Sleep(1 * time.Second) <--- ???????
	for _, contact := range shortList.Contacts {
		if contact.ID.Equals(kademlia.RoutingTable.me.ID) {
			result := kademlia.GetData(dataHash)
			return result
		} else {
			result, err := kademlia.Network.getAtTarget(dataKey, &contact)
			if err != nil {
				fmt.Println("Error getting data:", err.Error())
			} else if result != nil {
				return result
			}
		}
	}
	return nil
	//store the value in the closest node that isn't the correct node
}

func (kademlia *Kademlia) GetData(hash string) (data []byte) {
	result := kademlia.DataStore[hash]
	return result
}

func (kademlia *Kademlia) Store(data []byte) error {
	// get hash of data
	//dataHash = hex.EncodeToString(sha1.New().Sum(data))
	dataHash := hex.EncodeToString(hashData(data))
	dataKey := NewKademliaID(dataHash)

	// find closest nodes, Maybe use lookupcontact?
	shortlist := kademlia.LookupContact(dataKey, kademlia.RoutingTable.me)
	closest := shortlist.Contacts[0]

	if closest.ID.Equals(kademlia.RoutingTable.me.ID) {
		kademlia.StoreLocally(data, dataHash)
		return nil
	} else {
		kademlia.Network.storeAtTarget(data, &closest)
		return nil
	}
	//kademlia.Network.SendStoreMessage(data, closest)
	// store data in datastore
	// kademlia.DataStore[hash] = data
}

func (kademlia *Kademlia) StoreLocally(data []byte, dataHash string) {
	kademlia.DataStore[dataHash] = data
}
