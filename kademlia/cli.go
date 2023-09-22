package d7024e

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func CLI(kademlia *Kademlia) {
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
				execute(fieldedInput, ping, 1, "ping [node address]", kademlia)
			case "put":
				execute(fieldedInput, put, 2, "put [file]", kademlia)
			case "get":
				execute(fieldedInput, get, 2, "get [hash]", kademlia)
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
	newID := NewRandomKademliaID()
	contactInfo := NewContact(newID, input[1])
	err := kademlia.Network.SendPingMessage(&contactInfo)
	if err != nil {
		fmt.Printf("The ping was not successful. \n", err)
	}
}

func put(input []string, kademlia *Kademlia) {
	fmt.Printf("Your file was uploaded successfully! \n")
	err := kademlia.Network.SendStoreMessage([]byte(input[1]))
	if err != nil {
		fmt.Printf("your file was not uploaded successfully", err)
	} else {
		fmt.Printf("Your file was uploaded successfully! The id is: \n%s\n", input[1])
	}

}

func get(input []string, kademlia *Kademlia) {
	fmt.Printf("Your file was fetched succesfully! \n")
	//get file from hash here
	//return file to user here
}

func exit(input []string, kademlia *Kademlia) {
	fmt.Printf("Bye, bye little node! \n")
	//exit node here os.exit(0)
}
