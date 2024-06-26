package tcp

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

/* Here goes all the tcp functionalities */

type TCPCHAN[T string] struct {
    
    listenAddress string
    remoteAddress string

    InChan chan T
    OutChan chan T

    outBoundCon net.Conn
    listener net.Listener
}

func New[T string](lAddr, remoteAddr string) (*TCPCHAN[T], error) {
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
    go tcp.accLoop()
    go tcp.dialRemoteAndRead()

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

func (t *TCPCHAN[T]) accLoop() {

    defer func() {
        t.listener.Close()
    }()

    for {
        connect, err := t.listener.Accept()

        if err != nil {
            log.Println("Error accepting", err)
            return
        }

        log.Printf("sender connected %s", connect.RemoteAddr())

        go t.handleConnection(connect)

    }
}

func (t *TCPCHAN[T]) handleConnection(conn net.Conn) { 
    for {
        var messsage T

        if err := gob.NewDecoder(conn).Decode(&messsage); err != nil {
            log.Println(err)
            continue
        }

        t.OutChan <- messsage
    }
}

func (t *TCPCHAN[T]) dialRemoteAndRead() {
    
    connection, err := net.Dial("tcp", t.remoteAddress)

    if err != nil {

        log.Printf("dial error (%s)", err)
        time.Sleep(time.Second * 3)
        t.dialRemoteAndRead()
    }
    t.outBoundCon = connection
}

