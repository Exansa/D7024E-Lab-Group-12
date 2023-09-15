package d7024e

import (
	"fmt"
	"net"
)

type Network struct {
	node Kademlia
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

func Listen(ip string, port int) {
	// TODO
	listen, err := net.Listen("udp", ip+":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return err
	}
	// Close the listener when the application closes.
	defer listen.Close()
	fmt.Println("Listening on %s:%v", ip, port)
	for {
		// Listen for an incoming connection.
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return err
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func sendMessage(msg *RPC) {
	//possible encoding
	conn, err := net.Dial("udp")

	if err != nil {
		fmt.Printf("Error: %s", err)
		return
	}
	//send msg
	//
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
