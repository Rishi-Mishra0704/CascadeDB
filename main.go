package main

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"github.com/Rishi-Mishra0704/distributed-cas/p2p"
	"github.com/Rishi-Mishra0704/distributed-cas/store"
)

func handleRead(key string, casStore *store.Store) {
	// Read data from the store
	reader, err := casStore.Read(key)
	if err != nil {
		log.Fatalf("Error reading data from store: %v\n", err)
	}

	// Print the data read from the store
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, reader)
	if err != nil {
		log.Fatalf("Error copying data: %v\n", err)
	}
	fmt.Printf("Data read from store: %s\n", buf.String())
}

func handleWrite(key string, data []byte, casStore *store.Store) {
	// Write data to the store
	err := casStore.WriteStream(key, bytes.NewReader(data))
	if err != nil {
		log.Fatalf("Error writing data to store: %v\n", err)
	}
	fmt.Printf("Data written successfully to store\n")
}

func OnPeer(peer p2p.Peer) error {
	peer.Close()
	return nil
}

func main() {
	// Initialize your store with appropriate options
	storeOpts := store.StoreOpts{
		PathTransformFunc: store.CasPathTransformFunc,
	}
	casStore := store.NewStore(storeOpts)

	// Example usage of store methods
	key := "example_key"
	data := []byte("example data")

	// Check if key exists in the store
	exists := casStore.Has(key)
	fmt.Printf("Key exists: %v\n", exists)

	// Write data to the store
	handleWrite(key, data, casStore)

	// Check if key exists after writing
	exists = casStore.Has(key)
	fmt.Printf("Key exists after writing: %v\n", exists)

	// Read data from the store
	handleRead(key, casStore)

	// Create and start the TCP transport
	tcpOpts := p2p.TcpTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	// Start consuming messages from the transport
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	// Keep the main goroutine alive
	select {}
}
