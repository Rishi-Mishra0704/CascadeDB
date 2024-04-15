package server

import (
	"fmt"
	"log"
	"sync"

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
		case msg := <-s.Transport.Consume():
			fmt.Println(msg)

		case <-s.QuitCh:
			return
		}

	}
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
