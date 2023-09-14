package d7024e

type Network struct {
	Kademlia Kademlia
}

func Listen(ip string, port int) {
	// TODO
}

func sendMessage(msg *RPC) {
	//possible encoding
	//connect to other node
	//send msg
	//
	//check for errors between all
}

func (network *Network) SendPingMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.msgType = "PING"
	//newMsg.sender = *contact
	newMsg.data.PING = "Ping!"
	sendMessage(newMsg)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	newMsg := new(RPC)
	newMsg.msgType = "STORE"
	//newMsg.sender = *contact
	//newMsg.data.NODE = data
}

func (network *Network) SendFindDataMessage(hash string) {
	newMsg := new(RPC)
	newMsg.msgType = "STORE"
	//newMsg.sender = *contact
	newMsg.data.VALUE = hash
}

func (network *Network) SendStoreMessage(contact *Contact, data []byte) []byte {
	hash := hashData(data)
	newMsg := new(RPC)
	newMsg.msgType = "STORE"
	//newMsg.sender = *contact
	newMsg.data.STORE = data
	sendMessage(newMsg)
	return hash
}
