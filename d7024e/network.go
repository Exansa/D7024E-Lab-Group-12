package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net"
	"time"
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
		data := network.Kademlia.LookupData(msg.Data.VALUE)
		if data != nil {
			network.SendFoundDataMessage(data, &msg.Sender)
		}

	case FOUND_VALUE:
		//TODO:
		network.msgChan <- *msg
		network.dataChan <- msg.Data.STORE

	case ERR:
		network.msgChan <- *msg
		fmt.Println("Error:", msg.Data.ERR)
	case GET:
		data := network.Kademlia.GetData(msg.Data.HASH)
		if data != nil {
			network.SendFoundDataMessage(data, &msg.Sender)
		}
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

func (network *Network) storeAtTarget(data []byte, target *Contact) error {
	fmt.Print("Storing at target\n")
	hash := hex.EncodeToString(hashData(data))
	network.SendStoreMessage(data, hash, target)
	res := <-network.msgChan

	if res.Type != STORED || !res.Sender.ID.Equals(target.ID) {
		return fmt.Errorf("storeValue failed")
	}

	return nil
}

func (network *Network) getAtTarget(hash *KademliaID, target *Contact) ([]byte, error) {
	fmt.Print("Getting at target", target.Address, "\n")
	network.SendGetMessage(hash, target)
	res := <-network.msgChan

	if res.Type != FOUND_VALUE || !res.Sender.ID.Equals(target.ID) {
		return nil, fmt.Errorf("getValue failed")
	}

	return res.Data.VALUE, nil
}

func (network *Network) ping(contact *Contact) error {
	//TODO: Add timeout

	// Timeout after 5 seconds
	for i := 0; i < 10; i++ {
		fmt.Printf("Pinging %s\n", contact.Address)
		network.SendPingMessage(contact)
		fmt.Printf("Sent ping to %s\n", contact.Address)

		select {
		case res := <-network.msgChan:
			if res.Type == PONG && res.Sender.ID.Equals(contact.ID) {
				return nil
			} else {
				fmt.Printf("Ping failed!\n")
				return fmt.Errorf("ping failed")
			}
		case <-time.After(1 * time.Second):
			fmt.Printf("Ping timed out!\n")
			continue
		}
	}

	return fmt.Errorf("host unreachable")
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
