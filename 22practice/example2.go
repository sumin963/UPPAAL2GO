package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

var x_start time.Time // for clock x

var ch_reset chan bool = make(chan bool)        // for channel reset
var ch_reset_update chan bool = make(chan bool) // for update
var committed int = 0                           // for committed location
var lock_committed sync.Mutex

//var ch_progress chan bool // for update
var eps time.Duration = time.Millisecond
var lock_global sync.Mutex

func wait_for_committed() {
	for committed > 0 {
		runtime.Gosched()
	}
}

func Test() {
	var delay time.Duration
L1:

	fmt.Println("Test L1 starts")
	wait_for_committed()

	fmt.Println(x_start.Unix())

	<-time.After(2*time.Second - time.Since(x_start)) // x==2
	goto L1_1
L1_1:
	wait_for_committed()

	fmt.Println("Test L1_1 starts")

	delay = time.Duration(rand.Int31n(3000)) // for x>=2, how delaY?? 3000?
	<-time.After(delay * time.Millisecond)
	goto L1_2

L1_2:
	fmt.Println("Test L1_2 starts")

	wait_for_committed()
	ch_reset <- true

	lock_global.Lock()
	// update code here.
	lock_global.Unlock()

	ch_reset_update <- true // Observer's turn for update

	goto L1

}
func Observer() {
L1:
	fmt.Println("Observer L1 starts")
	<-ch_reset
	lock_committed.Lock()
	committed++
	lock_committed.Unlock()
	<-ch_reset_update
	lock_global.Lock()
	// update code here.
	lock_global.Unlock()

	goto L2

L2: // committed location
	fmt.Println("Observer L2 starts")
	lock_global.Lock()
	x_start = time.Now()
	lock_global.Unlock()

	lock_committed.Lock()
	committed--
	lock_committed.Unlock()
	goto L1

}
func main() {
	// wg := new(sync.WaitGroup)        // for waitGroup
	x_start = time.Now()

	// wg.Add(2)
	go Test()
	go Observer()

	time.Sleep(30 * time.Second)
	// wg.Wait()
}
