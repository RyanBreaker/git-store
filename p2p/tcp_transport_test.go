package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	listenAddress := ":4000"
	tr := NewTCPTransport(TCPTransportOpts{
		ListenAddr: ":4000",
	})

	assert.Equal(t, listenAddress, tr.ListenAddr)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
