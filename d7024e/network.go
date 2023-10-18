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
	msgBuffer    chan RPC
	dataChan     chan []byte
	lookupBuffer LookupBuffer
}

func NewNetwork(kademlia *Kademlia) *Network {
	network := Network{}
	network.Kademlia = kademlia
	network.msgBuffer = make(chan RPC, 100)
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

	// Add contact
	network.Kademlia.RoutingTable.AddContact(msg.Sender)

	// switch case for different message types
	switch msg.Type {
	case PING:
		network.SendPongMessage(&msg.Sender)

	case PONG:
		//TODO:
		network.msgBuffer <- *msg

	case STORE:
		// store data using kademlia func store
		network.Kademlia.StoreLocally(msg.Data.STORE, msg.Data.HASH)
		network.SendStoredMessage(&msg.Sender)

	case STORED:
		network.msgBuffer <- *msg

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
		network.msgBuffer <- *msg

	case FIND_VALUE:
		// based on hash, find data using kademlia func lookupdata
		data := network.Kademlia.LookupData(msg.Data.VALUE)
		if data != nil {
			network.SendFoundDataMessage(data, &msg.Sender)
		}

	case FOUND_VALUE:
		//TODO:
		network.msgBuffer <- *msg
		network.dataChan <- msg.Data.STORE

	case ERR:
		network.msgBuffer <- *msg
	case GET:
		data := network.Kademlia.GetData(msg.Data.HASH)
		if data != nil {
			network.SendFoundDataMessage(data, &msg.Sender)
		}
	default:
		fmt.Println("Message type not found")
	}
}

func (network *Network) awaitAndValidate(mt msgType, sender *KademliaID, timeout int) (RPC, error) {
	for i := 0; i < timeout*100; i++ {
		select {
		case res := <-network.msgBuffer:
			if (res.Type != mt && res.Type != ERR) || !res.Sender.ID.Equals(sender) {
				network.msgBuffer <- res
			} else {
				return res, nil
			}
		case <-time.After(time.Duration(timeout) * time.Second):
			return RPC{}, fmt.Errorf("timed out waiting for response from %s", sender)
		}
	}

	network.msgBuffer = make(chan RPC, 100)

	return RPC{}, fmt.Errorf("too many tries waiting on response from %s", sender)
}

func (network *Network) findNode(target *KademliaID, sender *Contact) (ContactCandidates, error) {
	network.SendFindContactMessage(target, sender)
	res, err := network.awaitAndValidate(FOUND_NODE, sender.ID, 5)

	if err != nil {
		return ContactCandidates{}, err
	}

	if res.Type == ERR {
		return ContactCandidates{}, fmt.Errorf("findNode failed: %s", res.Data.ERR)
	}

	//TODO: Dead statement
	if res.Type != FOUND_NODE || !res.Sender.ID.Equals(sender.ID) {
		return ContactCandidates{}, fmt.Errorf("findNode failed")
	}

	return res.Data.NODES, nil
}

func (network *Network) storeAtTarget(data []byte, target *Contact) error {
	hash := hex.EncodeToString(hashData(data))
	network.SendStoreMessage(data, hash, target)
	res, err := network.awaitAndValidate(STORED, target.ID, 5)

	if err != nil {
		return err
	}

	//TODO: Dead statement
	if res.Type != STORED || !res.Sender.ID.Equals(target.ID) {
		return fmt.Errorf("storeValue failed")
	}

	return nil
}

func (network *Network) getAtTarget(hash *KademliaID, target *Contact) ([]byte, error) {
	network.SendGetMessage(hash, target)
	res, err := network.awaitAndValidate(FOUND_VALUE, target.ID, 5)

	if err != nil {
		return nil, err
	}

	//TODO: Dead statement
	if res.Type != FOUND_VALUE || !res.Sender.ID.Equals(target.ID) {
		return nil, fmt.Errorf("getValue failed")
	}

	return res.Data.VALUE, nil
}

func (network *Network) ping(timeout int, contact *Contact) error {

	// Timeout after n retries
	for i := 0; i < timeout; i++ {
		network.SendPingMessage(contact)

		select {
		case res := <-network.msgBuffer:
			if res.Type == PONG && res.Sender.ID.Equals(contact.ID) {
				return nil
			} else {
				network.msgBuffer <- res
				continue
			}
		case <-time.After(1 * time.Second):
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
	for _, process := range buffer.lookups {
		if process.Sender.ID.Equals(msg.Sender.ID) && process.Data.NODE.Equals(&msg.Data.NODE) {
			return true
		}
	}
	return false
}

func (buffer *LookupBuffer) Remove(msg RPC) {
	for i, process := range buffer.lookups {
		if process.Sender.ID.Equals(msg.Sender.ID) && process.Data.NODE.Equals(&msg.Data.NODE) {
			buffer.lookups[i] = buffer.lookups[buffer.Len()-1] // Replace index with last elem
			buffer.lookups = buffer.lookups[:buffer.Len()-1]   // Truncate slice
		}
	}
}

func (buffer LookupBuffer) Len() int {
	return len(buffer.lookups)
}
