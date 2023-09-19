package d7024e

type Kademlia struct {
	ID           int               //id
	IP           string            //ip
	PORT         int               //port
	DataStore    map[string][]byte //data storage
	Bootstrap    bool              //bootstrap eller inte
	RoutingTable *RoutingTable     //routingtable
	//bucket?
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
	shortlist := kademlia.RoutingTable.FindClosestContacts(target.ID, 3)
	closest := shortlist[0]
	for _, contact := range shortlist {
		// Contact remote node and add the returned closest nodes to closestNodes
		// You need to implement the logic for contacting remote nodes and getting their closest nodes
		// You can use the kademlia.RoutingTable to find the closest nodes to a target

		if contact.ID.CalcDistance(target.ID).Less(closest.ID.CalcDistance(target.ID)) {
			closest = contact
		}
	}
	//async FIND_NODE RPC to the closest nodes in shortlist
	//wait for responses
	//k closest nodes per call and put them into shortlist
	//sort and remove duplicates in shortlist and cut at k
	//repeat until no closer nodes are found
	/*

			probed = Contact []

			for {
				lastClosest := closest

				for _, contact := range shortlist {
					if contact in probed
						continue
					probed.append(contact)

					res, err = kademlia.network.SendFindContactMessage(target, shortlist[n].ID, me)
					res.sort()

					addToShortList(res, shortlist)
				}
				shortlist.append(res)
				shortlist.sort()
				shortlist = shortlist[:k]

				if(closest == lastClosest){ // Add failbreak
					break
				}
			}
			return closest

		func addToShortList(res []Contact, shortlist []Contact) {
			shortlist.append(res)

		}
	*/
	return closest
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
