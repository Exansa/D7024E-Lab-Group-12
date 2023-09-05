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

		input = strings.TrimSpace(input)

		fmt.Printf("You entered: %s\n", input)
	}
}
