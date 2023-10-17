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
	PING  string            `json:"ping"`  // ping message
	STORE []byte            `json:"store"` // store message
	HASH  string            `json:"hash"`  // hash message
	NODE  KademliaID        `json:"node"`  // node message
	NODES ContactCandidates `json:"nodes"` // nodes message
	VALUE []byte            `json:"value"` // value message
	ERR   string            `json:"err"`   // error message
}

type msgType string

const (
	PING        msgType = "PING"
	PONG        msgType = "PONG" //Ack
	STORE       msgType = "STORE"
	STORED      msgType = "STORED"
	FIND_NODE   msgType = "FIND_NODE"
	FOUND_NODE  msgType = "FOUND_NODE"
	FIND_VALUE  msgType = "FIND_VALUE"
	FOUND_VALUE msgType = "FOUND_VALUE"
	ERR         msgType = "ERR" //Duplicate message
	GET         msgType = "GET"
	EXIT        msgType = "EXIT"
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
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Ping!"
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendPongMessage(contact *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = PONG
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Pong!"
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendFindContactMessage(target *KademliaID, receiver *Contact) (ContactCandidates, error) { //trasig
	newMsg := new(RPC)
	newMsg.Type = FIND_NODE
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.NODE = *target
	sendMessage(newMsg)
	return ContactCandidates{}, nil
}

func (network *Network) SendFoundContactMessage(contacts ContactCandidates, receiver *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = FOUND_NODE
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.NODES = contacts
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendFindDataMessage(hash []byte) error {
	newMsg := new(RPC)
	newMsg.Type = FIND_NODE
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Data.VALUE = hash
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendFoundDataMessage(data []byte, receiver *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = FOUND_VALUE
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.VALUE = data
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendStoreMessage(data []byte, hash string, receiver *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = STORE
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.STORE = data
	newMsg.Data.HASH = hash
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendStoredMessage(receiver *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = STORED
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendError(receiver *Contact, errStr string) {
	newMsg := new(RPC)
	newMsg.Type = ERR
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.ERR = errStr
	sendMessage(newMsg)
}

func (network *Network) SendGetMessage(hash *KademliaID, receiver *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = GET
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	newMsg.Data.HASH = hash.String()
	sendMessage(newMsg)
	return nil
}

func (network *Network) SendExitMessage(receiver *Contact) error {
	newMsg := new(RPC)
	newMsg.Type = EXIT
	newMsg.Sender = *network.Kademlia.RoutingTable.me
	newMsg.Receiver = *receiver
	sendMessage(newMsg)
	return nil
}
