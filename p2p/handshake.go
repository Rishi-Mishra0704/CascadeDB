package p2p

import "errors"

/*
	ErrInvalidHandshakeError is returned when a connection
	between local and remote nodes is not established
*/
var ErrInvalidHandshakeError = errors.New("invalid handshake")

// HandShakeFunc is a function that is called when a new connection is established
type HandShakeFunc func(Peer) error

func NOPHandShakeFunc(Peer) error { return nil }
