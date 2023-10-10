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

func (network *Network) Listen() {
	udpAddr, err := net.ResolveUDPAddr("udp", network.Kademlia.ADDRESS)
	checkError(err)

	fmt.Println("Listening on", udpAddr.IP)

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
	fmt.Println("Handling request from", msg.Sender.Address)
	fmt.Println("Message type:", msg.Type)

	// switch case for different message types
	switch msg.Type {
	case PING:
		network.Kademlia.RoutingTable.AddContact(msg.Sender)
		fmt.Printf("Received ping from %s\n", msg.Sender.ID.String())
		network.SendPongMessage(&msg.Sender)

	case PONG:
		//TODO:
		network.msgChan <- *msg

	case STORE:
		// store data using kademlia func store
		network.Kademlia.StoreValue(msg.Data.STORE, msg.Data.HASH)

	case FIND_NODE:
		// send closest nodes using kademlia func lookupcontact
		fmt.Println("Received find node message from", msg.Sender.ID.String())
		contacts := network.Kademlia.LookupContact(&msg.Data.NODE)
		fmt.Println("Found contacts:", contacts)
		fmt.Println("Sending found node message to", msg.Sender.ID.String())
		network.SendFoundContactMessage(contacts, &msg.Sender)

	case FOUND_NODE:
		//TODO:
		network.msgChan <- *msg

	case FIND_VALUE:
		// based on hash, find data using kademlia func lookupdata
		data := network.Kademlia.LookupData(msg.Data.VALUE)
		if data != "" {
			network.SendFoundDataMessage(data, &msg.Sender)
		}

	case FOUND_VALUE:
		//TODO:
		network.msgChan <- *msg
		network.dataChan <- msg.Data.STORE

	default:
		fmt.Println("Message type not found")
	}
}

func (network *Network) findNode(target *KademliaID, sender *Contact) (ContactCandidates, error) {
	network.SendFindContactMessage(target, sender)
	fmt.Println("Sent find node message to", sender.ID.String())
	res := <-network.msgChan
	fmt.Println("Received response from", res.Sender.ID.String())
	fmt.Println("Message:", res.Data.NODES)

	if res.Type != FOUND_NODE || !res.Sender.ID.Equals(sender.ID) {
		return ContactCandidates{}, fmt.Errorf("findNode failed")
	}

	return res.Data.NODES, nil
}

func (network *Network) ping(contact *Contact) error {
	//TODO: Add timeout

	fmt.Println("Sending ping to", contact.ID.String())

	network.SendPingMessage(contact)
	res := <-network.msgChan

	fmt.Println("Received response from", res.Sender.ID.String())

	if res.Type == PONG && res.Sender.ID.Equals(contact.ID) {
		// Add contact to routing table
		network.Kademlia.RoutingTable.AddContact(res.Sender)
		return nil
	} else {
		return fmt.Errorf("ping failed")
	}
}
