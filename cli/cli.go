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
		fmt.Printf("Enter a string: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input", err)
			return
		}

		trimedInput := strings.TrimSpace(input)
		fieldedInput := strings.Fields(trimedInput)
		switch fieldedInput[0] {
		case "put":
			if len(fieldedInput) == 2 {
				put(fieldedInput)
			} else {
				fmt.Printf("Wrong arguments! correct command is: put [file] \n")
			}
		case "get":
			if len(fieldedInput) == 2 {
				put(fieldedInput)
			} else {
				fmt.Printf("Wrong arguments! correct command is: get [hash] \n")
			}

		case "exit":
			if len(fieldedInput) == 1 {
				exit(fieldedInput)
			} else {
				fmt.Printf("Wrong arguments! correct command is: exit \n")
			}
		case "help":
			fmt.Printf("here are the different commands")
		}

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
