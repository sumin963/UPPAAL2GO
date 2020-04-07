package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	go func() {
		time.Sleep(time.Second * 30)
		c <- 2
	}()
	select {
	case <-c:
		fmt.Println("c")
	default:
		fmt.Println("d")

	}
}
