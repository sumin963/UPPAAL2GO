package main

import (
	"fmt"
	"reflect"
)

func produce(ch chan<- string, i int) {
	for j := 0; j < 5; j++ {
		ch <- fmt.Sprint(i*10 + j)
	}
	close(ch)
}
func main() {
	var sendCh = make(chan int)

	var increaseInt = func(c chan int) {
		for i := 0; i < 8; i++ {
			c <- i
		}
		close(c)
	}

	go increaseInt(sendCh)

	var selectCase = make([]reflect.SelectCase, 1)
	selectCase[0].Dir = reflect.SelectRecv
	selectCase[0].Chan = reflect.ValueOf(sendCh)

	counter := 0
	for counter < 1 {
		// use of Select() method
		chosen, recv, recvOk := reflect.Select(selectCase)
		if recvOk {
			fmt.Println(chosen, recv.Int(), recvOk)

		} else {
			counter++
		}
	}
}

// func main() {
// 	numChans := 4

// 	//I keep the channels in this slice, and want to "loop" over them in the select statemnt
// 	var chans = []chan string{}

// 	for i := 0; i < numChans; i++ {
// 		ch := make(chan string)
// 		chans = append(chans, ch)
// 		go produce(ch, i+1)
// 	}

// 	cases := make([]reflect.SelectCase, len(chans))
// 	for i, ch := range chans {
// 		cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ch)}
// 	}

// 	remaining := len(cases)
// 	for remaining > 0 {
// 		chosen, value, ok := reflect.Select(cases)
// 		if !ok {
// 			// The chosen channel has been closed, so zero out the channel to disable the case
// 			cases[chosen].Chan = reflect.ValueOf(nil)
// 			remaining -= 1
// 			continue
// 		}

// 		fmt.Printf("Read from channel %#v and received %s\n", chans[chosen], value.String())
// 	}
// }

// func main() {

// 	//chan1 := make(chan chan interface{})
// 	chan2 := make(chan interface{})
// 	chan3 := make(chan interface{})
// 	go aa(chan2, chan3)
// 	for i := 0; i < 10; i++ {
// 		select {
// 		case chan2 <- shortFunction():
// 			fmt.Println("1")
// 		case chan3 <- shortFunction():
// 			fmt.Println("2")
// 		}
// 		fmt.Println("ddd")
// 	}
// }
// func shortFunction() interface{} {
// 	defer fmt.Println("end short function")

// 	fmt.Println("start short function")
// 	return nil
// }
// func sFunction(chan1 chan chan interface{}) chan interface{} {
// 	defer fmt.Println("end chan chan function")

// 	fmt.Println("start chan chan function")

// 	responseChan := make(chan interface{})
// 	chan1 <- responseChan
// 	return <-chan1
// }
// func aa(chan1, chan2 chan interface{}) {
// 	for {
// 		select {
// 		case <-chan1:
// 			fmt.Println("11111")
// 		case <-chan2:
// 			fmt.Println("22222")
// 		}
// 	}
// }
