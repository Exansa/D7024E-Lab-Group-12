package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

const alpha = 3
const k = 20

type Kademlia struct {
	ID           int               //id
	IP           string            //ip
	PORT         int               //port
	DataStore    map[string][]byte //data storage
	Bootstrap    bool              //bootstrap eller inte
	RoutingTable *RoutingTable     //routingtable
	Network      *Network          //network
}

func (kademlia *Kademlia) LookupContact(target *Contact) Contact {

	shortlist := ContactCandidates{}
	shortlist.contacts = kademlia.RoutingTable.FindClosestContacts(target.ID, 3)
	closest := shortlist.contacts[0]
	probed := make(map[string]bool)
	probed[closest.ID.String()] = true
	queue := make(chan ContactCandidates, k)

	for {
		lastClosest := closest

		for _, contact := range shortlist.contacts {
			// Skip if already probed to prevent dupes
			if _, ok := probed[contact.ID.String()]; ok {
				continue
			}
			probed[contact.ID.String()] = true

			//async FIND_NODE RPC to the closest nodes in shortlist
			go func(contact *Contact) {
				res, err := kademlia.Network.SendFindContactMessage(target, contact)

				if err != nil {
					fmt.Println("Error listening:", err.Error())
					return // If it fails to reply, it won't be added to the shortlist
				}

				queue <- res

			}(&contact)

		}

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

	return data
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
