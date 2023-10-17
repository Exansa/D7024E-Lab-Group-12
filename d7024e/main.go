package main

import (
	"log"
	"net"
	"os"
)

func main() {
	ip := GetLocalIP()
	node := NewKademlia(ip.String() + ":8000")
	go node.initNode()
	CLI(os.Stdin, node)
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
