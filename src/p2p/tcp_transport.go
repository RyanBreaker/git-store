package p2p

import (
	"fmt"
	"io"
	"net"
)

// TCPPeer represents the remote node over a TCP-established connection.
type TCPPeer struct {
	// conn is the connection to the peer
	conn net.Conn

	// if we dial a connection => output == true
	// if we accept a connection => outbound == false
	outbound bool
}

func NewTCPPeer(conn net.Conn, outbound bool) *TCPPeer {
	return &TCPPeer{
		conn:     conn,
		outbound: outbound,
	}
}

func (p *TCPPeer) Close() error {
	return p.conn.Close()
}

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
	OnPeer        func(Peer) error
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener
	rpcCh    chan RPC
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
		rpcCh:            make(chan RPC),
	}
}

// Consume implements the Transport interface, returning a read-only channel for reading the incoming
// messages received from another peer on the network.
func (t *TCPTransport) Consume() <-chan RPC {
	return t.rpcCh
}

func (t *TCPTransport) ListenAndAccept() error {
	var err error

	t.listener, err = net.Listen("tcp", t.ListenAddr)
	if err != nil {
		return err
	}

	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) acceptLoop() error {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("TCP Listener accept error: %v\n", err)
			continue
		}

		fmt.Printf("New incoming connection: %+v\n", conn)

		go t.handleConn(conn)
	}
}

type Temp struct{}

func (t *TCPTransport) handleConn(conn net.Conn) {
	var err error

	defer func() {
		fmt.Printf("Dropping peer connection: %+v\n", conn)
		conn.Close()
	}()

	peer := NewTCPPeer(conn, false)

	if err := t.HandshakeFunc(peer); err != nil {
		return
	}

	if t.OnPeer != nil {
		if err = t.OnPeer(peer); err != nil {
			return
		}
	}

	// Read loop
	rpc := RPC{}
	for {
		if err = t.Decoder.Decode(conn, &rpc); err != nil {
			if err == io.EOF {
				fmt.Printf("TCP connection closed by peer: %+v\n", peer)
				break
			}
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}

		rpc.From = conn.RemoteAddr()
		t.rpcCh <- rpc
	}
}
