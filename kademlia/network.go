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
	MsgType  string  `json:"msgType"`  // message type
	Data     msgData `json:"data"`     // message data
}

type msgData struct {
	PING  string
	STORE []byte
	NODE  KademliaID
	VALUE string
}

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
	switch msg.MsgType {
	case "PING":
		// send pong
		if msg.Data.PING == "Ping!" {
			network.SendPongMessage(&msg.Sender)
		} else if msg.Data.PING == "Pong!" {
			// add sender to kademlia routing table
			network.Kademlia.RoutingTable.AddContact(msg.Sender)
		}
	case "STORE":
		// store data using kademlia func store
		network.Kademlia.Store(msg.Data.STORE)
	case "FIND_NODE":
		// send closest nodes
	case "FIND_VALUE":
		// based on hash, find data using kademlia func lookupdata

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

func sendMessage(msg *RPC) {

	conn, err := net.Dial("udp", msg.Receiver.Address)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	encodedMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	//send msg
	conn.Write(encodedMsg)
	//check for errors between all
	err = conn.Close()
	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.MsgType = "PING"
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Ping!"
	sendMessage(newMsg)
}

func (network *Network) SendPongMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.MsgType = "PING"
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.PING = "Pong!"
	sendMessage(newMsg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.MsgType = "STORE"
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	//newMsg.data.NODE = data
}

func (network *Network) SendFindDataMessage(hash string) {
	newMsg := new(RPC)
	newMsg.MsgType = "STORE"
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Data.VALUE = hash
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) []byte {
	hash := hashData(data)
	newMsg := new(RPC)
	newMsg.MsgType = "STORE"
	newMsg.Sender = network.Kademlia.RoutingTable.me
	newMsg.Receiver = *contact
	newMsg.Data.STORE = data
	sendMessage(newMsg)
	return hash
}
