package p2p

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTCPTransport(t *testing.T) {
	listenAddress := ":4000"
	tr := NewTCPTransport(listenAddress)

	assert.Equal(t, listenAddress, tr.listenAddress)

	// Server
	assert.Nil(t, tr.ListenAndAccept())
}
