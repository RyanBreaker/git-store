package p2p

// A Peer represents a remote node.
type Peer interface {
}

// Transport is anything that handles communications between nodes on the network.
// This can be TCP, UDP, websockets, etc.
type Transport interface {
	ListenAndAccept() error
}
