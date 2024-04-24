package main

import (
	"fmt"
	"log"
	"time"

	"github.com/MikelGV/tcp/tcp"
)

/* Here goes all the tcp functionalities */

func main() {
    localChannel, err := tcp.New[string]("localhost:3000", "localhost:8080")

    if err != nil {
        log.Fatal(err)
        return
    }

    go func() {
        time.Sleep(5 * time.Second)
        localChannel.InChan <- "test 1"
    }()

    remoteChannel, err := tcp.New[string]("localhost:8080", "localhost:3000" )

    if err != nil {
        log.Fatal(err)
        return
    }

    message := <-remoteChannel.OutChan

    fmt.Println("received from channel over wire:", message)
}
