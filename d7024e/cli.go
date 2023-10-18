package main

import (
	"bufio"
	"fmt"
	"io"
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
			case "awaitmsg":
				fmt.Printf("Awaiting message\n")
				res := <-kademlia.Network.msgChan
				fmt.Printf("Message received from %s with type %s\n", res.Sender.Address, res.Type)
			case "sendmsg":
				fmt.Printf("Sending message\n")
				contact := NewContact(NewKademliaID(bootstrapIDString), fieldedInput[1])
				kademlia.Network.SendError(&contact, "test error")
				fmt.Printf("Message sent\n")
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
	newID := NewKademliaID(bootstrapIDString)
	contactInfo := NewContact(newID, input[1])
	err := kademlia.Network.ping(&contactInfo)
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
	fmt.Printf("Your file was fetched succesfully! \n")
	res := kademlia.LookupData([]byte(input[1]))
	if res != nil {
		fmt.Println("The ping was not successful. \n", string(res))
	}
	//get file from hash here
	//return file to user here
}

func exit(input []string, kademlia *Kademlia) {
	fmt.Printf("Bye, bye little node! \n")
	kademlia.Network.SendExitMessage(kademlia.RoutingTable.me)
	//exit node here os.exit(0)
}
