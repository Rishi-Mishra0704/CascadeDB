package main

import (
	"bytes"
	"log"
	"time"

	"github.com/Rishi-Mishra0704/distributed-cas/p2p"
	"github.com/Rishi-Mishra0704/distributed-cas/server"
)

func makeServer(listenAddr string, nodes ...string) *server.FileServer {

	tcpTransportOpts := p2p.TcpTransportOpts{
		ListenAddr:    listenAddr,
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := server.FileServerOpts{
		StorageRoot:       listenAddr + "_network",
		PathTransformFunc: server.CasPathTransformFunc,
		Transport:         tcpTransport,
		BootstrapNodes:    nodes,
	}
	s := server.NewFileServer(fileServerOpts)
	tcpTransport.OnPeer = s.OnPeer
	return s
}

func main() {
	s1 := makeServer(":3000", "")
	s2 := makeServer(":4000", ":3000")
	go func() {
		log.Fatal(s1.Start())
	}()
	time.Sleep(4 * time.Second)
	go s2.Start()
	time.Sleep(4 * time.Second)
	data := bytes.NewReader([]byte("some very big data file!!!"))
	s2.StoreData("myprivatebigdata", data)

	select {}
}
