package main

import (
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
	Train := func(id int) {
		x_now := time.Now()
		x := time.Since(x_now)
		goto aa
	aa:
	}
	Gate := func(id int) {
		goto aa
	aa:
	}
}
func when(guard bool, channel chan bool) chan bool {
	if !guard {
		return nil
	}
	return channel
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
