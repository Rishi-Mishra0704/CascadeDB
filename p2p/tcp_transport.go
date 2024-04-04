package p2p

import (
	"fmt"
	"log"
	"net"
)

// TCPPeer represents the remote node over a TCP established connection
type TCPPeer struct {

	// conn is the underlying connection of the peer
	conn net.Conn

	// if we dial and retrieve a conn => outbound == true
	// if we accept and retrieve a conn => outbound == false
	outbound bool
}

type TcpTransportOpts struct {
	ListenAddr    string
	HandShakeFunc HandShakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TcpTransport struct {
	TcpTransportOpts
	listener net.Listener
	RpcCh    chan RPC
}

func NewTcpPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

func NewTCPTransport(opts TcpTransportOpts) *TcpTransport {
	return &TcpTransport{
		TcpTransportOpts: opts,
		RpcCh:            make(chan RPC),
	}
}

// Consume implements the transport interface, which will return read-only channel
// for reading the message recieved from another peer in the network
func (t *TcpTransport) Consume() <-chan RPC {
	return t.RpcCh
}

func (t *TcpTransport) ListenAndAccept() error {

	var err error
	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		log.Fatal(err)
	}

	go t.startAcceptLoop()

	return nil
}

func (t *TcpTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}
		fmt.Printf("New incoming connection %v\n", conn)

		go t.handleConn(conn)

	}
}

func (t *TcpTransport) handleConn(conn net.Conn) {

	var err error
	defer func() {
		fmt.Printf("Dropping peer connection%s", err)
		conn.Close()
	}()
	peer := NewTcpPeer(conn, true)

	if err := t.HandShakeFunc(peer); err != nil {

		return
	}

	// Read loop
	rpc := RPC{}
	for {
		err := t.Decoder.Decode(conn, &rpc)
		if err != nil {

			return
		}
		rpc.From = conn.RemoteAddr()
		t.RpcCh <- rpc
	}

}
