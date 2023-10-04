package d7024e

import (
	"encoding/json"
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

func sendMessage(msg *RPC) {
	// Dial sender
	conn, err := net.Dial("udp", msg.Receiver.Address)
	checkError(err)
	defer conn.Close()

	// Marshall msg
	jsonMsg, err := json.Marshal(msg)
	checkError(err)

	// Send msg
	_, err = conn.Write(jsonMsg)
	checkError(err)

}

func (network *Network) SendPingMessage(contact *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = PING
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Ping!"
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendPongMessage(contact *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = PING
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Pong!"
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendFindContactMessage(target *Contact, receiver *Contact) (ContactCandidates, error) { //trasig
	newMsg := new(RPC)
	newMsg.Type = FIND_NODE
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.NODE = *target
	sendMessage(newMsg)
	return ContactCandidates{}, nil
}

func (network *Network) SendFindDataMessage(hash string) error {
	newMsg := new(RPC)
	newMsg.Type = FIND_NODE
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Data.VALUE = hash
	sendMessage(newMsg)
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
		sendMessage(newMsg)
	}
	return nil
}
