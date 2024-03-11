package p2p

import "net"

/*

Message hold any arbitrary data that is being sent over each transport
between two nodes in a network


*/
type Message struct {
	From    net.Addr
	Payload []byte
}
