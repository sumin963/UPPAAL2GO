package main

import (
	"fmt"
	"math/rand"
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
			case x1_1 > time.Duration(5*time.Second):
				goto Alarm
			case x1_1 >= time.Duration(3*time.Second) && x1_1 <= time.Duration(5*time.Second):
				time.Sleep(time.Second * time.Duration(rand.Intn(2)))
				ev <- struct{}{}
				goto End
			}
		}

	End: // End1 Location
		fmt.Println("P1_End Location")
		goto Start

	Alarm:
		fmt.Println("P1_Alarm")

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
