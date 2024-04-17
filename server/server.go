package server

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/Rishi-Mishra0704/distributed-cas/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
	BootstrapNodes    []string
}

type FileServer struct {
	FileServerOpts

	Peerlock sync.Mutex
	Peers    map[string]p2p.Peer
	Store    *Store
	QuitCh   chan struct{}
}

func NewFileServer(opts FileServerOpts) *FileServer {
	storeOpts := StoreOpts{
		Root:              opts.StorageRoot,
		PathTransformFunc: opts.PathTransformFunc,
	}

	return &FileServer{
		FileServerOpts: opts,
		Store:          NewStore(storeOpts),
		QuitCh:         make(chan struct{}),
		Peers:          make(map[string]p2p.Peer),
	}
}

type Message struct {
	Payload any
}

type MessageStoreFile struct {
	Key  string
	Size int64
}

func (s *FileServer) broadcast(msg *Message) error {
	peers := []io.Writer{}

	for _, peer := range s.Peers {
		peers = append(peers, peer)
	}
	mw := io.MultiWriter(peers...)
	return gob.NewEncoder(mw).Encode(msg)

}

func (s *FileServer) StoreData(key string, r io.Reader) error {
	buf := new(bytes.Buffer)
	tee := io.TeeReader(r, buf)
	size, err := s.Store.Write(key, tee)
	if err != nil {
		return err
	}

	msgBuf := new(bytes.Buffer)
	msg := Message{
		Payload: MessageStoreFile{
			Key:  key,
			Size: size,
		},
	}

	if err := gob.NewEncoder(msgBuf).Encode(msg); err != nil {
		return err
	}

	for _, peer := range s.Peers {
		if err := peer.Send(msgBuf.Bytes()); err != nil {
			return err
		}

	}
	time.Sleep(2 * time.Second)
	for _, peer := range s.Peers {
		n, err := io.Copy(peer, buf)
		if err != nil {
			return err
		}
		fmt.Printf("recieved and written bytes to disk: %v\n", n)
	}

	return nil

}

func (s *FileServer) Stop() {
	close(s.QuitCh)

}

func (s *FileServer) OnPeer(p p2p.Peer) error {
	s.Peerlock.Lock()
	defer s.Peerlock.Unlock()

	s.Peers[p.RemoteAddr().String()] = p

	log.Printf("connected with remote: %s", p.RemoteAddr())
	return nil
}

func (s *FileServer) Loop() {
	defer func() {
		fmt.Println("Closing the file server due to quit action")
		s.Transport.Close()
	}()

	for {
		select {
		case rpc := <-s.Transport.Consume():
			var msg Message
			if err := gob.NewDecoder(bytes.NewReader(rpc.Payload)).Decode(&msg); err != nil {
				log.Fatal("error decoding: ", err)
			}
			if err := s.HandleMessage(rpc.From, &msg); err != nil {
				log.Println("error handling message: ", err)
				return
			}

		case <-s.QuitCh:
			return
		}

	}
}

func (s *FileServer) HandleMessage(from string, msg *Message) error {
	switch v := msg.Payload.(type) {
	case MessageStoreFile:
		return s.HandleMessageStoreFile(from, v)
	}
	return nil
}
func (s *FileServer) HandleMessageStoreFile(from string, msg MessageStoreFile) error {
	peer, ok := s.Peers[from]
	if !ok {
		return fmt.Errorf("peer (%s) not found in peer map", from)
	}
	if _, err := s.Store.Write(msg.Key, io.LimitReader(peer, msg.Size)); err != nil {
		return nil
	}
	peer.(*p2p.TCPPeer).Wg.Done()
	return nil
}
func (s *FileServer) bootstrapNetwork() error {
	for _, addr := range s.BootstrapNodes {
		if len(addr) == 0 {
			continue
		}
		go func(addr string) {
			fmt.Println("attempting to connect with: ", addr)
			if err := s.Transport.Dial(addr); err != nil {
				log.Println("dial err", err)

			}
		}(addr)
	}
	return nil
}

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.bootstrapNetwork()

	s.Loop()
	return nil
}

func init() {
	gob.Register(MessageStoreFile{})
}
