package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

type Network struct {
	Kademlia Kademlia
}

type RPC struct {
	Sender   Contact `json:"sender"`   // sender Contact
	Receiver Contact `json:"receiver"` // receiver Contact
	Type     msgType `json:"msgType"`  // message type
	Data     msgData `json:"data"`     // message data
}

type msgData struct {
	PING  string
	STORE []byte
	HASH  string
	NODE  Contact
	VALUE string
}

type msgType string

const (
	PING       msgType = "PING"
	STORE      msgType = "STORE"
	FIND_NODE  msgType = "FIND_NODE"
	FIND_VALUE msgType = "FIND_VALUE"
)

func (network *Network) handleRequest(conn net.Conn) {
	// read incoming message and decode it
	packet := make([]byte, 1024)
	_, err := conn.Read(packet)
	if err != nil {
		fmt.Printf("Error: %s", err)
		conn.Close()
		return
	}
	var msg RPC
	err = json.Unmarshal(packet, &msg)
	// check for errors
	if err != nil {
		fmt.Printf("Error: %s", err)
		conn.Close()
		return
	}
	// switch case for different message types
	switch msg.Type {
	case PING:
		// send pong
		if msg.Data.PING == "Ping!" {
			network.SendPongMessage(&msg.Sender)
		} else if msg.Data.PING == "Pong!" {
			// add sender to kademlia routing table
			network.Kademlia.RoutingTable.AddContact(msg.Sender)
		}
	case STORE:
		// store data using kademlia func store
		network.Kademlia.StoreValue(msg.Data.STORE, msg.Data.HASH)

	case FIND_NODE:
		// send closest nodes
	case FIND_VALUE:
		// based on hash, find data using kademlia func lookupdata
		network.Kademlia.LookupData(msg.Data.VALUE)

	default:
		fmt.Println("Message type not found")
	}
	// send response
	// check for errors
	// close connection
	conn.Close()

}

func (network *Network) Listen(address string) {
	// TODO
	listen, err := net.Listen("udp", address) //kolla på att använda net.addr
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	// Close the listener when the application closes.
	defer listen.Close()
	fmt.Println("Listening on ", address)
	for {
		// Listen for an incoming connection.
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}
		// Handle connections in a new goroutine.
		go network.handleRequest(conn)
	}
}

func sendMessage(msg *RPC) error {

	conn, err := net.Dial("udp", msg.Receiver.Address)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	//send msg
	conn.Write(encodedMsg)
	//check for errors between all
	err = conn.Close()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return nil
}

func (network *Network) SendPingMessage(contact *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = PING
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Ping!"
	err := sendMessage(newMsg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return nil
}

func (network *Network) SendPongMessage(contact *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = PING
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Pong!"
	err := sendMessage(newMsg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return nil
}

func (network *Network) SendFindContactMessage(target *Contact, receiver *Contact) (ContactCandidates, error) { //trasig
	newMsg := new(RPC)
	newMsg.Type = FIND_NODE
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.NODE = *target
	err := sendMessage(newMsg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return ContactCandidates{}, err
	}
	return ContactCandidates{}, nil
}

func (network *Network) SendFindDataMessage(hash string) error {
	newMsg := new(RPC)
	newMsg.Type = FIND_NODE
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Data.VALUE = hash
	err := sendMessage(newMsg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	return nil
}

func (network *Network) SendStoreMessage(data []byte) error {
	newMsg := new(RPC)
	newMsg.Type = STORE
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Data.STORE = data
	self, receiver, dataHash := network.Kademlia.Store(newMsg.Data.STORE)
	if !self {
		newMsg.Receiver = receiver
		newMsg.Data.HASH = dataHash
		err := sendMessage(newMsg)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return err
		}
	}
	return nil
}
