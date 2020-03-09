package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type chan_t struct{}
var wg sync.WaitGroup

func main(){
	ev1 := make(chan chan_t)
	r := make(chan chan_t)
	ctx, cancel := context.WithCancel(context.Background())

	wg.Add(1)
	go p(ctx,ev1,r,cancel)
	go c(ctx,ev1,r,cancel)

	fmt.Println("main go")
	wg.Wait()
	fmt.Println("main done")
}
func p (ctx context.Context,ev1 chan chan_t,r chan chan_t,cancel func()) {
	time.Sleep(time.Second*1)
	fmt.Println("p sync")
	ev1<-chan_t{}
	for {
		select {
		case <-r:
			time.Sleep(time.Second*1)
			fmt.Println("p sync")
			ev1<-chan_t{}
		case <-ctx.Done():
			fmt.Println("p down")
			wg.Done()
		}
	}
}
func c (ctx context.Context,ev1 chan chan_t,r chan chan_t,cancel func()) {
	for {
		select {
		case <-time.After(time.Second*2):
			<-ev1
			r<-chan_t{}
			fmt.Println("C ack")
		case <-time.After(time.Second*1):
			fmt.Println("Timeout")
		    cancel()
			return
		case <-ctx.Done():
		    wg.Done()
		}
	}
}