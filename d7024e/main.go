package d7024e

import (
	"log"
	"net"
	"os"
)

func main() {
	ip := GetSwarmIP()
	node := NewKademlia(ip + ":8000")
	go node.initNode()
	CLI(os.Stdin, node)
}

// Get the local IP address for the swarm service
func GetSwarmIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.String()[:3] == "10." {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}
