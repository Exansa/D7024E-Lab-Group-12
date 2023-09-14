package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		trimedInput := strings.TrimSpace(input)
		fieldedInput := strings.Fields(trimedInput)
		if len(fieldedInput) > 0 {
			switch fieldedInput[0] {
			case "put":
				execute(fieldedInput, put, 2, "put [file]")
			case "get":
				execute(fieldedInput, get, 2, "get [hash]")
			case "exit":
				execute(fieldedInput, exit, 1, "exit")
			case "help":
				fmt.Printf("here are the different commands")
			default:
				fmt.Printf("Invalid command.\n")
			}
		}

	}
}

func execute(inp []string, exec func([]string), inpLen int, corrStr string) {
	if len(inp) == inpLen {
		exec(inp)
	} else {
		fmt.Printf("Invalid argument\nCorrect format: %s\n\n", corrStr)
	}
}

func put(input []string) {
	fmt.Printf("Your file was uploaded succesfully! \n")
	//store value here
	//return hash here
}

func get(input []string) {
	fmt.Printf("Your file was fetched succesfully! \n")
	//get file from hash here
	//return file to user here
}

func exit(input []string) {
	fmt.Printf("Bye, bye little node! \n")
	//exit node here
}
