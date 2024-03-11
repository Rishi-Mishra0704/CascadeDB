package p2p

import (
	"encoding/gob"
	"io"
)

type Decoder interface {
	Decode(io.Reader, *RPC) error
}

type GobDecoder struct {
}

func (doc GobDecoder) Decode(r io.Reader, m *RPC) error {
	return gob.NewDecoder(r).Decode(m)
}

type DefaultDecoder struct {
}

func (doc DefaultDecoder) Decode(r io.Reader, m *RPC) error {
	buf := make([]byte, 1028)
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	m.Payload = buf[:n]
	return nil

}
