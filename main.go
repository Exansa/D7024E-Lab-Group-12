package main

import (
	"d7024e"
	"log"
	"net"
	"os"
)

func main() {
	ip := GetLocalIP()
	node := d7024e.NewKademlia(ip.String() + ":8000")
	node.initNode()
	node.CLI(os.Stdin)
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
