package d7024e

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func CLI(stdin io.Reader, kademlia *Kademlia) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		trimmedInput := strings.TrimSpace(input)
		fieldedInput := strings.Fields(trimmedInput)
		if len(fieldedInput) > 0 {
			switch fieldedInput[0] {
			case "ping":
				execute(fieldedInput, ping, 2, "ping [node address]", kademlia)
			case "put":
				execute(fieldedInput, put, 2, "put [file]", kademlia)
			case "get":
				execute(fieldedInput, get, 2, "get [hash]", kademlia)
			case "ip":
				fmt.Printf("Your ip is: %s\n", kademlia.ADDRESS)
			case "id":
				fmt.Printf("Your id is: %s\n", kademlia.ID)
			case "info":
				fmt.Println("--------------------NODE INFO---------------------")
				fmt.Printf("ID: %s\n", kademlia.ID)
				fmt.Printf("IP: %s\n", kademlia.ADDRESS)
				fmt.Printf("Bucket size: %d\n", len(kademlia.RoutingTable.buckets))
				fmt.Printf("Datastore size: %d\n", len(kademlia.DataStore))
				fmt.Println("--------------------------------------------------")
			case "exit":
				execute(fieldedInput, exit, 1, "exit", kademlia)
			case "help":
				fmt.Printf("here are the different commands")
			default:
				fmt.Printf("Invalid command.\n")
			}
		}

	}
}

func execute(inp []string, exec func([]string, *Kademlia), inpLen int, corrStr string, kademlia *Kademlia) {
	if len(inp) == inpLen {
		exec(inp, kademlia)
	} else {
		fmt.Printf("Invalid argument\nCorrect format: %s\n\n", corrStr)
	}
}

func ping(input []string, kademlia *Kademlia) {
	// Check if the input is a valid address
	address := input[1]
	_, err := net.ResolveIPAddr("udp", address)
	if err != nil {
		fmt.Printf("Invalid address: %s\n", address)
		return
	}

	newID := NewKademliaID(bootstrapIDString)
	contactInfo := NewContact(newID, input[1])
	err = kademlia.Network.ping(5, &contactInfo)
	if err != nil {
		fmt.Println("The ping was not successful. \n", err)
	} else {
		fmt.Printf("The ping was successful! \n")
	}
}

func put(input []string, kademlia *Kademlia) {
	//fmt.Printf("Your file was uploaded successfully! \n")
	err := kademlia.Store([]byte(input[1]))
	if err != nil {
		fmt.Println("your file was not uploaded successfully", err)
	} else {
		fmt.Printf("Your file was uploaded successfully! The id is: \n%s\n", input[1])
	}

}

func get(input []string, kademlia *Kademlia) {
	res := kademlia.LookupData([]byte(input[1]))
	if res != nil {
		fmt.Println("Success! Found value: \n", string(res))
	}
	//get file from hash here
	//return file to user here
}

func exit(input []string, kademlia *Kademlia) {
	fmt.Printf("Bye, bye little node! \n")
	if input[0] == "test" {

	} else {
		os.Exit(0)
	}
}
