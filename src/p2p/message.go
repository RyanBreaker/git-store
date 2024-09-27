package p2p

import "net"

// An RPC holds any arbitrary data that is being sent between two nodes on the network.
type RPC struct {
	From    net.Addr
	Payload []byte
}
