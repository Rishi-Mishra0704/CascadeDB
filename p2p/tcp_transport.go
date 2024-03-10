package p2p

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type TcpTransport struct {
	listenAddress string
	listener      net.Listener
	mu            sync.RWMutex
	peers         map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TcpTransport {
	return &TcpTransport{
		listenAddress: listenAddr,
	}
}

func (t *TcpTransport) ListenAndAccept() error {

	var err error
	t.listener, err = net.Listen("tcp", t.listenAddress)
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

		go t.handleConn(conn)

	}
}

func (t *TcpTransport) handleConn(conn net.Conn) {
	fmt.Printf("New Incoming Connection %+v\n", conn)
}
