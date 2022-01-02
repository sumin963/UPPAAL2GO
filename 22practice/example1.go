package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var y_start time.Time // for clock x

var ch_press chan bool = make(chan bool)        // for channel reset
var ch_press_update chan bool = make(chan bool) // for update ordering

var committed int = 0 // for committed location
var lock_committed sync.Mutex

//var ch_progress chan bool // for update
var eps time.Duration = time.Millisecond
var lock_global sync.Mutex

func wait_for_committed() {
	for committed > 0 {
		runtime.Gosched()
	}
}

func Lamp() {
	var delay time.Duration
L_off:
	fmt.Println("Lamp's off starts")
	fmt.Println(y_start.Unix())
	wait_for_committed()

	<-ch_press
	<-ch_press_update
	lock_global.Lock()
	y_start = time.Now()
	lock_global.Unlock()
	goto L_low

L_low:
	select {
	case <-ch_press:
		<-ch_press_update
		lock_global.Lock()
		// update code here.
		lock_global.Unlock()
		goto L_bright
	case <-time.After(5*time.Second - time.Since(y_start)): // y==5
		goto L_low_1
	}

L_low_1:
	wait_for_committed()
	fmt.Println("At Lamp's low, 5 seconds passed")

	delay = time.Duration(rand.Int31n(3000)) // for x>=2, how delaY?? 3000?
	<-time.After(delay * time.Millisecond)

	<-ch_press
	<-ch_press_update
	lock_global.Lock()
	// update code here.
	lock_global.Unlock()
	goto L_off

L_bright:
	wait_for_committed()
	fmt.Println("At Lamp's bright starts")

	delay = time.Duration(rand.Int31n(3000)) // for x>=2, how delaY?? 3000?
	<-time.After(delay * time.Millisecond)

	<-ch_press
	<-ch_press_update
	lock_global.Lock()
	// update code here.
	lock_global.Unlock()
	goto L_off
}

func User() {
	var delay time.Duration
L_idle:

	fmt.Println("User's idle starts")
	wait_for_committed()

	delay = time.Duration(rand.Int31n(10000)) // for x>=2, how delaY?? 3000?
	<-time.After(delay * time.Millisecond)

	ch_press <- true
	fmt.Println("Sync occurs after", time.Since(y_start))

	lock_global.Lock()
	// update code here.
	lock_global.Unlock()
	ch_press_update <- true

	goto L_idle
}
func main() {
	// wg := new(sync.WaitGroup)        // for waitGroup
	y_start = time.Now()

	// wg.Add(2)
	go Lamp()
	go User()

	time.Sleep(50 * time.Second)
	// wg.Wait()
}
