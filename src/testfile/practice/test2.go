package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ev := make(chan interface{})
	P1 := func() {
		x1 := time.Now()

	start:
		fmt.Println("P1 start location")
		time.Sleep(time.Second * 1) //Guard

		for {
			x1_1 := time.Since(x1)
			switch {
			case int(x1_1)/1000000000 > 3: //Invariant
				goto Error

			default:
				select {
				case ev <- struct{}{}:
					fmt.Println("P1 End location")
					x1 = time.Now()
					goto start
				default:
				}
			}
		}
	Error:
		fmt.Println("Error")
		return
	}
	P2 := func() {
		//x2 := time.Now()
		for {
			fmt.Println("P2 start location")
			time.Sleep(time.Second * 1) //Guard
			select {
			case <-ev:
				fmt.Println("P2 End location")
				//x2 = time.Now()
			case <-time.After(time.Second * 2): //Invariant
				goto Error

			}
		}
	Error:
		fmt.Println("Error")
		return
		//wg.Done()
	}

	wg.Add(1)
	go P1()
	go P2()
	wg.Wait()

}
