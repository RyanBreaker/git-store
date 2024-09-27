package p2p

import (
	"fmt"
	"io"
	"net"
	"sync"
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

type TCPTransportOpts struct {
	ListenAddr    string
	HandshakeFunc HandshakeFunc
	Decoder       Decoder
}

type TCPTransport struct {
	TCPTransportOpts
	listener net.Listener

	peerLock sync.RWMutex
	peers    map[net.Addr]Peer
}

func NewTCPTransport(opts TCPTransportOpts) *TCPTransport {
	return &TCPTransport{
		TCPTransportOpts: opts,
	}
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
	peer := NewTCPPeer(conn, true)

	if err := t.HandshakeFunc(peer); err != nil {
		conn.Close()
		fmt.Printf("TCP handshake error: %s\n", err)
		return
	}

	// Read loop
	msg := &RPC{}
	for {
		if err := t.Decoder.Decode(conn, msg); err != nil {
			if err == io.EOF {
				fmt.Printf("TCP connection closed by peer %+v\n", peer)
				break
			}
			fmt.Printf("TCP Error: %s\n", err)
			continue
		}

		msg.From = conn.RemoteAddr()

		fmt.Printf("RPC: %+v\n", msg)
	}
}
