package tcp

/* Here goes all the tcp functionalities */

type TCPCHAN[T any] struct {
    InChan chan T
    OutChan chan T
}

