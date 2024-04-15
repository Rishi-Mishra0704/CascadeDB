package p2p

import (
	"errors"
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

// RemoteAddr implements the peer interface
// which will return the remote address of its underlying connection
func (p *TCPPeer) RemoteAddr() net.Addr {
	return p.conn.RemoteAddr()
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

// Send implements the peer interface, which will send the message to the remote peer
func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
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

	log.Printf("Tcp Transport listening on port%s\n", t.ListenAddr)

	return nil
}

// Dial implements the transport interface
func (t *TcpTransport) Dial(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil
	}
	go t.handleConn(conn, true)
	return nil

}

// Close implements the transport interface, which will close the connection
func (t *TcpTransport) Close() error {
	return t.listener.Close()
}

func (t *TcpTransport) startAcceptLoop() {
	for {
		conn, err := t.listener.Accept()

		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			fmt.Printf("TCP accept error %s\n", err)
		}
		fmt.Printf("New incoming connection %v\n", conn)

		go t.handleConn(conn, false)

	}
}

func (t *TcpTransport) handleConn(conn net.Conn, outbound bool) {

	var err error
	defer func() {
		fmt.Printf("Dropping peer connection%s", err)
		conn.Close()
	}()
	peer := NewTcpPeer(conn, outbound)

	if err := t.HandShakeFunc(peer); err != nil {

		return
	}

	if t.OnPeer != nil {
		if err := t.OnPeer(peer); err != nil {
			return
		}
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
