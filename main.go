package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

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
func handlePut(conn net.Conn, casStore *store.Store) {
	// Read key from client
	scanner := bufio.NewScanner(conn)
	fmt.Fprint(conn, "Enter key: ")

	var key string

	// Read key synchronously
	for scanner.Scan() {
		key = strings.TrimSpace(scanner.Text())
		if key != "" {
			break
		}
	}

	fmt.Printf("Received key: %s\n", key)

	// Read data from client until a line break is encountered
	fmt.Fprint(conn, "Enter data: ")
	var dataBuf bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Skip empty lines
		}
		if line == "." {
			break // Stop reading if "." is entered
		}
		dataBuf.WriteString(line)
	}
	data := dataBuf.Bytes()

	// Write data to the store
	err := casStore.Put(key, bytes.NewReader(data))
	if err != nil {
		log.Printf("Error writing data to store: %v\n", err)
		return
	}

	fmt.Fprintln(conn, "Data written successfully to store")
}

func main() {
	// Initialize your store with appropriate options
	storeOpts := store.StoreOpts{
		PathTransformFunc: store.CasPathTransformFunc,
	}
	casStore := store.NewStore(storeOpts)

	// Create and start the TCP transport
	tcpOpts := p2p.TcpTransportOpts{
		ListenAddr:    ":3000",
		HandShakeFunc: p2p.NOPHandShakeFunc,
		Decoder:       p2p.DefaultDecoder{},
		OnPeer:        OnPeer,
		HandleConnCh:  make(chan p2p.HandleConnInfo),
	}
	tr := p2p.NewTCPTransport(tcpOpts)
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

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

	// Start consuming messages from the transport
	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()
	for {
		info := <-tr.HandleConnCh
		if info.Port == 3000 {
			go handlePut(info.Conn, casStore)
		}
	}
}
