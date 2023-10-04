package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

type Network struct {
	Kademlia *Kademlia
	msgChan  chan RPC
	dataChan chan []byte
}

func NewNetwork(kademlia *Kademlia) *Network {
	network := Network{}
	network.Kademlia = kademlia
	network.msgChan = make(chan RPC)
	network.dataChan = make(chan []byte)

	return &network
}

func (network *Network) Listen(addr string) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	defer conn.Close()

	for {
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		checkError(err)

		var msg RPC
		err = json.Unmarshal(buf[:n], &msg)
		checkError(err)

		go network.handleRequest(&msg)
	}
}

func (network *Network) handleRequest(msg *RPC) { // Server side

	// switch case for different message types
	switch msg.Type {
	case PING:
		network.SendPongMessage(&msg.Sender)

	case STORE:
		// store data using kademlia func store
		network.Kademlia.StoreValue(msg.Data.STORE, msg.Data.HASH)

	case FIND_NODE:
		// send closest nodes using kademlia func lookupcontact
		network.findNode(&msg.Data.NODE, &msg.Sender)

	case FOUND_NODE:
		//TODO:

	case FIND_VALUE:
		// based on hash, find data using kademlia func lookupdata
		data := network.Kademlia.LookupData(msg.Data.VALUE)
		if data != "" {
			network.SendFoundDataMessage(data, &msg.Sender)
		}

	case FOUND_VALUE:
		//TODO:

	default:
		fmt.Println("Message type not found")
	}
}

func (network *Network) findNode(target *KademliaID, sender *Contact) ContactCandidates {
	contacts := network.Kademlia.LookupContact(target)
	if contacts.Len() > 0 {
		network.SendFoundContactMessage(contacts, sender)
	}
	return contacts
}
