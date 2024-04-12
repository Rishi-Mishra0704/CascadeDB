package main

import (
	"log"
	"time"

	"github.com/Rishi-Mishra0704/distributed-cas/p2p"
	"github.com/Rishi-Mishra0704/distributed-cas/server"
)

func main() {

	tcpTransportOpts := p2p.TcpTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		// Todo OnPeer func
	}

	tcpTransport := p2p.NewTCPTransport(tcpTransportOpts)

	fileServerOpts := server.FileServerOpts{
		StorageRoot:       "3000_network",
		PathTransformFunc: server.CasPathTransformFunc,
		Transport:         tcpTransport,
	}
	fileServer := server.NewFileServer(fileServerOpts)
	go func() {
		time.Sleep(3 * time.Second)
		fileServer.Stop()
	}()
	if err := fileServer.Start(); err != nil {
		log.Fatalf("Error starting file server: %v", err)
	}

}
