package p2p

// A HandshakeFunc is...?
type HandshakeFunc func(any) error

func NOPHandshakeFunc(any) error {
	return nil
}
