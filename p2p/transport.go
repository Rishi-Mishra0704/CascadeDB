package p2p

// Peer is an interface that represents remote node.
type Peer interface {
	Close() error
}

/*
Trasnport is anything that controls the communication
between the nodes in the network
forms: tcp udp websockets ..etc
*/
type Transport interface {
	ListenAndAccept() error
	Consume() <-chan RPC
}
