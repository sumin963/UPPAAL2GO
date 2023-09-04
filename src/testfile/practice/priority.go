package main

import (
	"fmt"
	"time"
)

func main() {
	c := make(chan int)
	u := make(chan int)
	go func() {
		time.Sleep(time.Second * 2)
		c <- 2
	}()
	go func() {
		time.Sleep(time.Second * 2)
		u <- 2
	}()

	select {
	case <-c:
		fmt.Println("q")
	case <-u:
		select {
		case <-c:
			fmt.Println("c")
		default:
			fmt.Println("d")
		}
	}
}
