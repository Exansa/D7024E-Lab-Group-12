package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
)

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

func (network *Network) SendPongMessage(contact *Contact, conn *net.UDPConn) error {
	newMsg := new(RPC)
	newMsg.Type = PING
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Pong!"
	encodedMsg, err := json.Marshal(newMsg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return err
	}
	_, err = conn.Write(encodedMsg)
	// err := sendMessage(newMsg)
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
