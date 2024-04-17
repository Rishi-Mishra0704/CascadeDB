package p2p

/*

Message hold any arbitrary data that is being sent over each transport
between two nodes in a network


*/
type RPC struct {
	From    string
	Payload []byte
}
