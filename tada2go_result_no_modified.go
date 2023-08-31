package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// #define ctimemax {1,1,1,1}
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
	appr_chan := make([]chan bool, C.N)
	for i := range appr_chan {
		appr_chan[i] = make(chan bool)
	}
	stop_chan := make([]chan bool, C.N)
	for i := range stop_chan {
		stop_chan[i] = make(chan bool)
	}
	leave_chan := make([]chan bool, C.N)
	for i := range leave_chan {
		leave_chan[i] = make(chan bool)
	}
	go_chan := make([]chan bool, C.N)
	for i := range go_chan {
		go_chan[i] = make(chan bool)
	}
	Train := func(id C.id_t) {
		x_now := time.Now()
		x := time.Since(x_now)
		var id2_passage []string
		var id3_passage []string
		var id4_passage []string
	id0:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Safe", "location", "x", ":", x)
		select {
		case appr_chan[id] <- true:
			x_now = time.Now()
			goto id3
		}
	id1:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Stop", "location", "x", ":", x)
		select {
		case <-go_chan[id]:
			x_now = time.Now()
			goto id4
		}
	id2:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross_0", "location", "x", ":", x)
		id2_passage = []string{"x==3", "x>3", "x>5"}
		switch time_passage(id2_passage, x) {
		case 0:
		case 1:
			goto id2p
		case 2:
			goto id2pp
		case 3:
			goto exp
		}
		select {
		case <-time.After(time.Second*3 - x - eps):
			goto id2p
		}
	id3:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr_0", "location", "x", ":", x)
		id3_passage = []string{"x==10", "x>10", "x>20"}
		switch time_passage(id3_passage, x) {
		case 0:
		case 1:
			goto id3p
		case 2:
			goto id3pp
		case 3:
			goto exp
		}
		select {
		case <-time.After(time.Second*10 - x - eps):
			goto id3p
		case <-stop_chan[id]:
			goto id1
		}
	id4:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start_0", "location", "x", ":", x)
		id4_passage = []string{"x==7", "x>7", "x>15"}
		switch time_passage(id4_passage, x) {
		case 0:
		case 1:
			goto id4p
		case 2:
			goto id4pp
		case 3:
			goto exp
		}
		select {
		case <-time.After(time.Second*7 - x - eps):
			goto id4p
		}
	exp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "exp", "location", "x", ":", x)
		select {}
	id2p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross_1", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*3 - x):
			goto id2pp
		case leave_chan[id] <- true:
			goto id0
		}
	id2pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Cross_2", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*5 - x):
			goto exp
		case leave_chan[id] <- true:
			goto id0
		}
	id3p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr_1", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*10 - x):
			goto id3pp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			goto id2
		case <-stop_chan[id]:
			goto id1
		}
	id3pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Appr_2", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*20 - x):
			goto exp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			goto id2
		}
	id4p:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start_1", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*7 - x):
			goto id4pp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			goto id2
		}
	id4pp:
		x = time.Since(x_now)
		fmt.Println("Train", "template", "Start_2", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*15 - x):
			goto exp
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			goto id2
		}
	}
	Gate := func() {
		local_val := C.Gate{list: [C.N + 1]C.id_t{}, len: 0}
	id5:
		select {
		case stop_chan[C.tail(&local_val)] <- true:
			goto id6
		}
	id6:
		select {
		case <-time.After(time.Second * 40):
			C.enqueue(&local_val, e)
			goto id5

		case <-time.After(time.Second * 40):
			C.dequeue(&local_val)
			goto id7
		}
	id7:
		select {
		case when(local_val.len > 0, go_chan[C.front(&local_val)]) <- true:
			goto id6
		case <-time.After(time.Second * 40):
			C.enqueue(&local_val, e)
			goto id6
		}
	}
	go Train(0)
	go Gate()
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
