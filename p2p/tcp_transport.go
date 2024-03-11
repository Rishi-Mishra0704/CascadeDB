package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
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
}

type TcpTransport struct {
	TcpTransportOpts
	listener net.Listener

	mu    sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTcpPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func NewTCPTransport(opts TcpTransportOpts) *TcpTransport {
	return &TcpTransport{
		TcpTransportOpts: opts,
	}
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
	peer := NewTcpPeer(conn, true)
	if err := t.HandShakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error %s\n", err)
		return
	}
	// Read loop
	msg := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			fmt.Printf("TCP error %s\n", err)
			continue
		}
		msg.From = conn.RemoteAddr()
		fmt.Printf("Message : %+v\n", msg)
	}

}
