package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

//
import "C"

func main() {
	var closed1_chan chan bool
	var closed2_chan chan bool
	var activated1_chan chan bool
	var activated2_chan chan bool
	eps := time.Millisecond * 10
	Door := func() {
		x_now := time.Now()
		x := time.Since(x_now)
		goto id5
	id0:
		x = time.Since(x_now)
		fmt.Println("Door", "template", "wait", "location", "x", ":", x)
		select {
		case <-closed2_chan:
			x_now = time.Now()
			goto id4
		case closed1_chan <- true:
			goto id0
		}
	id1:
		x = time.Since(x_now)
		fmt.Println("Door", "template", "closing", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*6 - x - eps):
			x_now = time.Now()
			goto id2
		}
	id2:
		x = time.Since(x_now)
		fmt.Println("Door", "template", "closed", "location", "x", ":", x)
		select {
		case closed1_chan <- true:
			goto id2
		case <-time.After(time.Second*0 - x):
			goto id5
		}
	id3:
		x = time.Since(x_now)
		fmt.Println("Door", "template", "open", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*0 - x):
			x_now = time.Now()
			goto id1
		}
	id4:
		x = time.Since(x_now)
		fmt.Println("Door", "template", "opening", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*6 - x - eps):
			x_now = time.Now()
			goto id3
		}
	id5:
		x = time.Since(x_now)
		fmt.Println("Door", "template", "idle", "location", "x", ":", x)
		select {
		case closed1_chan <- true:
			goto id5
		case <-pushed_chan:
			C.activated = "true"
			goto id0
		}
	}
	User := func() {
		goto id7
	id6:
		select {
		case pushed_chan <- true:
			goto id7
		}
	id7:
		select {
		case <-time.After(time.Second * 40):
			C.w = "0"
			goto id6
		}
	}
	go Door()
	go Door()
	go User()
	go User()
	<-time.After(time.Second * 20)
}
func when(guard bool, channel chan bool) chan bool {
	if !guard {
		return nil
	}
	return channel
}
func when_guard(guard bool) <-chan time.Time {
	if !guard {
		return nil
	}
	return time.After(time.Second * 0)
}
func time_passage(time_passage []string, ctime time.Duration) int {
	for i, val := range time_passage {
		if strings.Contains(val, "==") {
			num, _ := strconv.Atoi(val[strings.Index(val, "==")+2:])
			if time.Second*time.Duration(num) > ctime {
				return i
			}
		} else if strings.Contains(val, "<") {
			num, _ := strconv.Atoi(val[strings.Index(val, "==")+1:])
			if time.Second*time.Duration(num) == ctime {
				return i
			}
		}
	}
	return len(time_passage)
}
