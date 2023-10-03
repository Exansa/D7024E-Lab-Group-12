package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type Network struct {
	Kademlia *Kademlia
	Client   *Client
	Server   *Server
}

type Client struct {
	Network *Network
	Out     chan RPC
	In      chan RPC
}

type Server struct {
	Network *Network
	Out     chan RPC
	In      chan RPC
}

func NewNetwork(kademlia *Kademlia) *Network {
	network := Network{}
	server := Server{}
	client := Client{}

	client.In = make(chan RPC)
	client.Out = make(chan RPC)
	client.Network = &network

	server.In = make(chan RPC)
	server.Out = make(chan RPC)
	server.Network = &network

	network.Kademlia = kademlia
	network.Client = &client
	network.Server = &server

	return &network
}

// Server SERVER SEEERE
func (network *Network) handleRequest(conn *net.UDPConn, addr *net.UDPAddr, buf []byte) { // Server side

	var msg RPC
	err := json.Unmarshal(buf, &msg)
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
		// if msg.Data.PING == "Ping!" {
		// 	network.SendPongMessage(&msg.Sender)
		// } else if msg.Data.PING == "Pong!" {
		// 	// add sender to kademlia routing table
		// 	network.Kademlia.RoutingTable.AddContact(msg.Sender)
		// }
		network.SendPongMessage(&msg.Sender, conn)

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

// Client send req
func (c *Client) SendRequest(address string, request []byte) ([]byte, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	_, err = conn.Write(request)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
		os.Exit(1)
	}
}

func (Client *Client) InitConnection(address string) {
	// TODO
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	checkError(err)

	dial, err := net.DialUDP("udp", nil, udpAddr)
	checkError(err)

	// Close the listener when the application closes.
	defer dial.Close()
	fmt.Println("Listening on ", address)
	for {
		// Listen for an incoming connection.
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			return
		}
		// Handle connections in a new goroutine.
		go Client.Network.handleRequest(dial, udpAddr)
	}
}

func (Server *Server) Start(addr string) {

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	checkError(err)

	conn, err := net.ListenUDP("udp", udpAddr)
	checkError(err)

	for {
		buf := make([]byte, 1024)
		n, addr, err := conn.ReadFromUDP(buf)
		checkError(err)

		go Server.Network.handleRequest(conn, addr, buf[:n])
	}
}
