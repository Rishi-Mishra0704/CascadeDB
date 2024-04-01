package p2p

import (
	"encoding/json"
	"net"
)

// CASCommand represents the command that will be sent over the network
type CASCommand struct {
	Operation string
	Key       string
	Payload   []byte
}

// SendCASCommand sends a CAS command message over the TCP connection
func SendCASCommand(conn net.Conn, command CASCommand) error {
	// Serialize command to JSON or other format
	commandBytes, err := json.Marshal(command)
	if err != nil {
		return err
	}
	// Send serialized command over TCP connection
	_, err = conn.Write(commandBytes)
	return err
}
