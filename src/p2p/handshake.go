package p2p

// A HandshakeFunc is...?
type HandshakeFunc func(Peer) error

func NOPHandshakeFunc(_ Peer) error {
	return nil
}
