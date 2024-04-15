package p2p

import "net"

// Peer is an interface that represents remote node.
type Peer interface {
	Send([]byte) error
	RemoteAddr() net.Addr
	Close() error
}

/*
Trasnport is anything that controls the communication
between the nodes in the network
forms: tcp udp websockets ..etc
*/
type Transport interface {
	Dial(string) error
	ListenAndAccept() error
	Consume() <-chan RPC
	Close() error
}
