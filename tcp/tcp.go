package tcp

import (
	"encoding/gob"
	"log"
	"net"
)

/* Here goes all the tcp functionalities */

type TCPCHAN[T any] struct {
    
    listenAddress string
    remoteAddress string

    InChan chan T
    OutChan chan T

    outBoundCon net.Conn
    listener net.Listener
}

func New[T any](lAddr, remoteAddr string) (*TCPCHAN[T], error) {
    tcp := &TCPCHAN[T]{
        listenAddress: lAddr,
        remoteAddress: remoteAddr,
        InChan: make(chan T, 10),
        OutChan: make(chan T, 10),
    }

    ln, err := net.Listen("tcp", lAddr)

    if err != nil {
        return nil, err
    }

    tcp.listener = ln

    go tcp.loop()

    return tcp, nil
}

func (t *TCPCHAN[T]) loop() {
    for {
        message := <-t.InChan
        log.Println("sending message", message)
        if err := gob.NewEncoder(t.outBoundCon).Encode(&message); err != nil {
            log.Println(err)
        }
    }
}

