package d7024e

import (
	"bytes"
	"fmt"
	"testing"
)

var buf bytes.Buffer

func TestCLIPing(t *testing.T) {
	fmt.Println("TestCLI")
	node := NewKademlia("10.0.0.2", true)

	buf.Write([]byte("ping 10.0.0.2"))
	CLI(&buf, node)
}

func TestCLIget(t *testing.T) {
	fmt.Println("TestCLI")
	node := NewKademlia("10.0.0.2", true)
	buf.Write([]byte("get"))
	CLI(&buf, node)
}

func TestCLIPut(t *testing.T) {
	fmt.Println("TestCLI")
	node := NewKademlia("10.0.0.2", true)
	buf.Write([]byte("put image.jpg"))
	CLI(&buf, node)
}

func TestCLIExit(t *testing.T) {
	fmt.Println("TestCLI")
	node := NewKademlia("10.0.0.2", true)
	buf.Write([]byte("exit"))
	CLI(&buf, node)
}

func TestCLIHelp(t *testing.T) {
	fmt.Println("TestCLI")
	node := NewKademlia("10.0.0.2", true)
	buf.Write([]byte("help"))
	CLI(&buf, node)
}
