package main

import (
	"d7024e"
	"os"
)

func start() {
	ip := os.Args[1]
	port := os.Args[2]
	address := ip + ":" + port
	node := d7024e.NewKademlia(address)
	node.initNode()
}
