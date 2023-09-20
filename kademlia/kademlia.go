package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
)

type Kademlia struct {
	ID           int               //id
	IP           string            //ip
	PORT         int               //port
	DataStore    map[string][]byte //data storage
	Bootstrap    bool              //bootstrap eller inte
	RoutingTable *RoutingTable     //routingtable
	Network      *Network          //network
}

/*
Backup

	func (kademlia *Kademlia) LookupContact(target *Contact) Contact {
		contacts := kademlia.RoutingTable.FindClosestContacts(target.ID, 3)
		if len(contacts) < 1 {
			return Contact{}
		}
		closest := contacts[0]
		for _, contact := range contacts {
			if closest.ID.Equals(target.ID) {
				return closest
			} else if contact.ID.CalcDistance(target.ID).Less(closest.ID.CalcDistance(target.ID)) {
				closest = contact
			}
		}
		return closest
	}
*/
func (kademlia *Kademlia) LookupContact(target *Contact) Contact {

	// for _, contact := range shortlist {
	// 	// Contact remote node and add the returned closest nodes to closestNodes
	// 	// You need to implement the logic for contacting remote nodes and getting their closest nodes
	// 	// You can use the kademlia.RoutingTable to find the closest nodes to a target

	// 	if contact.ID.CalcDistance(target.ID).Less(closest.ID.CalcDistance(target.ID)) {
	// 		closest = contact
	// 	}
	// }
	//async FIND_NODE RPC to the closest nodes in shortlist
	//wait for responses
	//k closest nodes per call and put them into shortlist
	//sort and remove duplicates in shortlist and cut at k
	//repeat until no closer nodes are found

	shortlist := kademlia.RoutingTable.FindClosestContacts(target.ID, 3)
	closest := shortlist[0]
	probed := []Contact{}

	for {
		lastClosest := closest

		for _, contact := range shortlist {
			if contains(probed, contact) {
				continue
			}

			probed = append(probed, contact)

			//TODO: async FIND_NODE RPC to the closest nodes in shortlist
			res, err := kademlia.Network.SendFindContactMessage(target, &contact)
			if err != nil {
				//TODO: handle error
			}
			res.Sort()
			closest = shortlist[0]

			//addToShortList(res, shortlist)

			//shortlist = append(shortlist, res)
			//shortlist.sort()

			//shortlist = shortlist[:k]

		}

		if closest == lastClosest { // Add failbreak
			break
		}
	}
	return closest

}

func contains[T comparable](s []T, e T) bool {
	for _, v := range s {
		if v == e {
			return true
		}
	}
	return false
}

func (kademlia *Kademlia) LookupData(hash string) {
	// similar to lookupcontact
}

func (kademlia *Kademlia) Store(data []byte) {
	// get hash of data
	dataHash := hex.EncodeToString(sha1.New().Sum(data))
	dataKey := NewKademliaID(dataHash)

	// find closest nodes, Maybe use lookupcontact?
	contacts := kademlia.RoutingTable.FindClosestContacts(dataKey, 3)
	closest := kademlia.RoutingTable.me
	// send store message to closest nodes
	for _, contact := range contacts {
		if contact.ID.CalcDistance(contact.ID).Less(closest.ID.CalcDistance(contact.ID)) {
			closest = contact
		}
	}
	if closest.ID.Equals(kademlia.RoutingTable.me.ID) {
		kademlia.StoreValue(data, dataHash)
	} else {
		kademlia.Network.SendStoreMessage(closest, data, dataHash)
	}
	//kademlia.Network.SendStoreMessage(data, closest)
	// store data in datastore
	// kademlia.DataStore[hash] = data
}

func (kademlia *Kademlia) StoreValue(data []byte, dataHash string) {
	kademlia.DataStore[dataHash] = data
}
