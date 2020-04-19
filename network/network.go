package network

import (
	"container/list"
)

/*
https://appliedgo.net/networking/

Data:
	rxMessage (shared queue)
	txMessage (shared queue)
	hosts

Functions:
	Listen (run as its own thread)
	Send (run as its own thread)

*/

type Message struct {

}

var (
	RxQueue = list.New()
	TxQueue = list.New()
)

func Foo() {
	RxQueue.PushBack(9)
}