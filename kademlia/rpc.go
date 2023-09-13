package d7024e

type RPC struct {
	sender  Contact
	msgType string
	data    msgData
}

type msgData struct {
	PING  string
	STORE []byte
	NODE  KademliaID
	VALUE string
}

func PING() {

}

func STORE() {

}

func FIND_NODE() {

}

func FIND_VALUE() {

}
