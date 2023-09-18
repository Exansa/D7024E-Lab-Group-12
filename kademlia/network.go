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
	sender   Contact
	receiver Contact
	msgType  string
	data     msgData
}

type msgData struct {
	PING  string
	STORE []byte
	NODE  KademliaID
	VALUE string
}

func (network *Network) Listen(ip string, port int) {
	// TODO
	listen, err := net.Listen("udp", ip+":"+port) //kolla på att använda net.addr
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}
	// Close the listener when the application closes.
	defer listen.Close()
	fmt.Println("Listening on %s:%v", ip, port)
	for {
		// Listen for an incoming connection.
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
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
	switch msg.msgType {
	case "PING":
		// send pong
	case "STORE":
		// store data
	case "FIND_NODE":
		// send closest nodes
	case "FIND_VALUE":
		// send value
	}
	// send response
	// check for errors
	// close connection
	conn.Close()

}

func sendMessage(msg *RPC) {

	conn, err := net.Dial("udp", msg.receiver.Address)
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
	newMsg.msgType = "PING"
	newMsg.sender = network.Kademlia.RoutingTable.me
	newMsg.receiver = *contact
	newMsg.data.PING = "Ping!"
	sendMessage(newMsg)
}

func (network *Network) SendPongMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.msgType = "PING"
	newMsg.sender = network.Kademlia.RoutingTable.me
	newMsg.receiver = *contact
	newMsg.data.PING = "Pong!"
	sendMessage(newMsg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.msgType = "STORE"
	newMsg.sender = network.Kademlia.RoutingTable.me
	newMsg.receiver = *contact
	//newMsg.data.NODE = data
}

func (network *Network) SendFindDataMessage(hash string) {
	newMsg := new(RPC)
	newMsg.msgType = "STORE"
	newMsg.sender = network.Kademlia.RoutingTable.me
	newMsg.data.VALUE = hash
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) []byte {
	hash := hashData(data)
	newMsg := new(RPC)
	newMsg.msgType = "STORE"
	newMsg.sender = network.Kademlia.RoutingTable.me
	newMsg.receiver = *contact
	newMsg.data.STORE = data
	sendMessage(newMsg)
	return hash
}
