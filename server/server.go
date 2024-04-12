package server

import (
	"fmt"

	"github.com/Rishi-Mishra0704/distributed-cas/p2p"
)

type FileServerOpts struct {
	ListenAddr        string
	StorageRoot       string
	PathTransformFunc PathTransformFunc
	Transport         p2p.Transport
}

type FileServer struct {
	FileServerOpts

	Store  *Store
	QuitCh chan struct{}
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
	}
}

func (s *FileServer) Stop() {
	close(s.QuitCh)

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

func (s *FileServer) Start() error {
	if err := s.Transport.ListenAndAccept(); err != nil {
		return err
	}

	s.Loop()
	return nil
}
