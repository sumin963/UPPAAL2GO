package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// Declaration
	ev := make(chan interface{}) // chan ev;

	// Template P1
	P1 := func() {
		// P1 Declaration
		x1 := time.Now() // clock x1;
		x1_1 := time.Since(x1)

	Start: // Start Location
		fmt.Println("P1_Start Location")
		time.Sleep(time.Second * 3)
		for {
			x1_1 = time.Since(x1)
			switch {
			case x1_1 >= time.Duration(3*time.Second):
				ev <- struct{}{}
				goto End
			}
		}

	End: // End1 Location
		fmt.Println("P1_End Location")
		x1 = time.Now()
		goto Start

	}

	// Template P2
	P2 := func() {
		// P2 Declaration

	Start: // Start Location
		fmt.Println("P2_Start Location")
		<-ev
		goto End

	End: // End1 Location
		fmt.Println("P2_End Location")
		goto Start
	}

	wg.Add(1)

	// System Declaration
	go P1()
	go P2()

	wg.Wait()
}
