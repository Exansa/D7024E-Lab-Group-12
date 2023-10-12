package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

type Network struct {
	Kademlia     *Kademlia
	msgChan      chan RPC
	dataChan     chan []byte
	lookupBuffer LookupBuffer
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

	// Add contact
	network.Kademlia.RoutingTable.AddContact(msg.Sender)

	// switch case for different message types
	switch msg.Type {
	case PING:
		network.SendPongMessage(&msg.Sender)

	case PONG:
		//TODO:
		network.msgChan <- *msg

	case STORE:
		// store data using kademlia func store
		network.Kademlia.StoreLocally(msg.Data.STORE, msg.Data.HASH)
		network.SendStoredMessage(&msg.Sender)

	case STORED:
		network.msgChan <- *msg

	case FIND_NODE:
		// send closest nodes using kademlia func lookupcontact
		if network.lookupBuffer.Has(*msg) {
			network.SendError(&msg.Sender, "DUPLICATE")
			return
		}

		network.lookupBuffer.Append(*msg)

		contacts := network.Kademlia.LookupContact(&msg.Data.NODE, &msg.Sender)
		network.SendFoundContactMessage(contacts, &msg.Sender)

		network.lookupBuffer.Remove(*msg)

	case FOUND_NODE:
		//TODO:
		network.msgChan <- *msg

	case FIND_VALUE:
		// based on hash, find data using kademlia func lookupdata
		can, data := network.Kademlia.LookupData(msg.Data.VALUE)
		if data != nil {
			network.SendFoundDataMessage(data, &msg.Sender)
		} else if can.Len() > 0 {
			network.SendFoundContactMessage(can, &msg.Sender)
		} else {
			network.SendError(&msg.Sender, "NOT_FOUND")
		}

	case FOUND_VALUE:
		//TODO:
		network.msgChan <- *msg
		network.dataChan <- msg.Data.STORE

	case ERR:
		network.msgChan <- *msg
		fmt.Println("Error:", msg.Data.ERR)

	default:
		fmt.Println("Message type not found")
	}
}

func (network *Network) findNode(target *KademliaID, sender *Contact) (ContactCandidates, error) {
	network.SendFindContactMessage(target, sender)
	res := <-network.msgChan

	if res.Type == ERR {
		return ContactCandidates{}, fmt.Errorf("findNode failed: %s", res.Data.ERR)
	}

	if res.Type != FOUND_NODE || !res.Sender.ID.Equals(sender.ID) {
		return ContactCandidates{}, fmt.Errorf("findNode failed")
	}

	return res.Data.NODES, nil
}

func (network *Network) findValue(target *string, sender *Contact) (ContactCandidates, []byte, error) {
	network.SendFindDataMessage(*target, sender)
	res := <-network.msgChan

	if !res.Sender.ID.Equals(sender.ID) {
		return ContactCandidates{}, nil, fmt.Errorf("findValue failed: Sender not equal")
	}

	switch res.Type {
	case FOUND_VALUE:
		return ContactCandidates{}, res.Data.STORE, nil
	case FOUND_NODE:
		return res.Data.NODES, nil, nil
	case ERR:
		return ContactCandidates{}, nil, fmt.Errorf("findValue failed: %s", res.Data.ERR)
	}

	return ContactCandidates{}, nil, fmt.Errorf("findValue failed")
}

func (network *Network) storeAtTarget(data []byte, target *Contact) error {
	network.SendStoreMessage(data, target)
	res := <-network.msgChan

	if res.Type != STORED || !res.Sender.ID.Equals(target.ID) {
		return fmt.Errorf("storeValue failed")
	}

	return nil
}

func (network *Network) ping(contact *Contact) error {
	//TODO: Add timeout

	network.SendPingMessage(contact)
	res := <-network.msgChan

	if res.Type == PONG && res.Sender.ID.Equals(contact.ID) {
		return nil
	} else {
		return fmt.Errorf("ping failed")
	}
}

type LookupBuffer struct {
	lookups []RPC
}

func (buffer *LookupBuffer) Append(msg RPC) {
	if msg.Type != FIND_NODE {
		return
	}

	buffer.lookups = append(buffer.lookups, msg)
}

func (buffer LookupBuffer) Has(msg RPC) bool {
	if msg.Type != FIND_NODE {
		return false
	}

	for _, process := range buffer.lookups {
		if process.Type != FIND_NODE {
			continue
		}

		if process.Sender.ID.Equals(msg.Sender.ID) && process.Data.NODE.Equals(&msg.Data.NODE) {
			return true
		}
	}
	return false
}

func (buffer *LookupBuffer) Remove(msg RPC) {
	for i, process := range buffer.lookups {
		if process.Sender == msg.Sender && process.Type == msg.Type && process.Data.NODE == msg.Data.NODE {
			buffer.lookups = append(buffer.lookups[:i], buffer.lookups[i+1:]...)
		}
	}
}

func (buffer LookupBuffer) Len() int {
	return len(buffer.lookups)
}
