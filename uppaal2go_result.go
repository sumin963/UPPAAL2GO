package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// #define N 6
// typedef  int id_t;
// typedef struct Train{
// }Train;
// typedef struct Gate{
//         id_t list[N+1];
//         int len;
// }Gate;
// void enqueue(Gate *Gate, id_t element)
// {
//         Gate->list[Gate-> len++] = element;
// }
// void dequeue(Gate *Gate )
// {
//         int i = 0;
//         Gate-> len -= 1;
//         while (i < Gate-> len)
//         {
//                 Gate->list[i] = Gate->list[i + 1];
//                 i++;
//         }
//         Gate->list[i] = 0;
// }
// id_t front(Gate *Gate )
// {
//    return Gate->list[0];
// }
// id_t tail(Gate *Gate )
// {
//    return Gate->list[Gate-> len - 1];
// }
//
import "C"

func main() {
	eps := time.Millisecond * 10
	appr_chan := make([]chan bool, "C.N")
	for i := range appr_chan {
		appr_chan[i] = make(chan bool)
	}
	stop_chan := make([]chan bool, "C.N")
	for i := range stop_chan {
		stop_chan[i] = make(chan bool)
	}
	leave_chan := make([]chan bool, "C.N")
	for i := range leave_chan {
		leave_chan[i] = make(chan bool)
	}
	go_chan := make([]chan bool, "C.")
	for i := range go_chan {
		go_chan[i] = make(chan bool)
	}
	Train := func(id C.id_t) {
		x_now := time.Now()
		x := time.Since(x_now)
	id0:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Safe", "location", "x", ":", x)
		select {
		case appr_chan[id] <- true:
			goto id3
		}
	id1:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Stop", "location", "x", ":", x)
		select {
		case <-go_chan[id]:
			goto id4
		}
	id2:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross", "location", "x", ":", x)
		select {
		case <-when_guard("x==3"):
			goto id2p
		}
	id3:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr", "location", "x", ":", x)
		select {
		case <-when_guard("x==10"):
			goto id3p
		case <-stop_chan[id]:
			goto id1
		}
	id4:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start", "location", "x", ":", x)
		select {
		case <-when_guard("x==7"):
			goto id4p
		}
	exp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "exp", "location", "x", ":", x)
		select {}
	id2p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "id2p", "location", "x", ":", x)
		select {
		case <-when_guard("x>3"):
			goto id2pp
		case leave_chan[id] <- true:
			goto id0
		}
	id2pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "id2pp", "location", "x", ":", x)
		select {
		case <-when_guard("x>5"):
			goto exp
		case leave_chan[id] <- true:
			goto id0
		}
	id3p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "id3p", "location", "x", ":", x)
		select {
		case <-when_guard("x>10"):
			goto id3pp
		case <-time.After(time.Second * 0):
			goto id2
		case <-stop_chan[id]:
			goto id1
		}
	id3pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "id3pp", "location", "x", ":", x)
		select {
		case <-when_guard("x>20"):
			goto exp
		case <-time.After(time.Second * 0):
			goto id2
		}
	id4p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "id4p", "location", "x", ":", x)
		select {
		case <-when_guard("x>7"):
			goto id4pp
		case <-time.After(time.Second * 0):
			goto id2
		}
	id4pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "id4pp", "location", "x", ":", x)
		select {
		case <-when_guard("x>15"):
			goto exp
		case <-time.After(time.Second * 0):
			goto id2
		}
	}
	Gate := func() {
		local_val := C.Gate{list: 0, len: 0}
	id5:
		select {
		case stop_chan[tail()] <- true:
			goto id6
		}
	id6:
		select {
		case <-time.After(time.Second * 40):
			goto id5
		case <-time.After(time.Second * 40):
			goto id7
		}
	id7:
		select {
		case when("len > 0", go_chan[front()]) <- true:
			goto id6
		case <-time.After(time.Second * 40):
			goto id6
		}
	}
	go Train()
	go Gate()
	<-time.After(time.Second * 20)
}
func when(guard bool, channel chan bool) chan bool {
	if !guard {
		return nil
	}
	return channel
}
func when_guard(guard bool) <-chan time.time {
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
