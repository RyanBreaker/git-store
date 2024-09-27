package main

import (
	"fmt"
	p2p2 "go-store/src/p2p"
	"log"
)

func OnPeer(p2p2.Peer) error {
	fmt.Println("Doing some logic...")
	return nil
}

func main() {
	tcpOpts := p2p2.TCPTransportOpts{
		ListenAddr:    ":3000",
		Decoder:       p2p2.DefaultDecoder{},
		HandshakeFunc: p2p2.NOPHandshakeFunc,
		OnPeer:        OnPeer,
	}
	tr := p2p2.NewTCPTransport(tcpOpts)

	go func() {
		for {
			msg := <-tr.Consume()
			fmt.Printf("%+v\n", msg)
		}
	}()

	fmt.Println("Starting up...")
	if err := tr.ListenAndAccept(); err != nil {
		log.Fatal(err)
	}

	select {}
}
