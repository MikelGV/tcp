package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/MikelGV/tcp/tcp"
)

/* Here goes all the tcp functionalities */

func main() {
    reader := bufio.NewReader(os.Stdin)
    line, err := reader.ReadString('\n')
    if err != nil {
        log.Fatal(err)
    }
    localChannel, err := tcp.New[string]("localhost:3000", "localhost:8080")

    if err != nil {
        log.Fatal(err)
        return
    }

    go func() {
        time.Sleep(5 * time.Second)
        localChannel.InChan <- line
    }()

    remoteChannel, err := tcp.New[string]("localhost:8080", "localhost:3000" )

    if err != nil {
        log.Fatal(err)
        return
    }

    message := <-remoteChannel.OutChan

    fmt.Println("received from channel over wire:", message)
}
