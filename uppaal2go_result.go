package main

import "time"

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
	appr_chan := make([]chan struct{})
	for i := range appr_chan {
		appr_chan[i] = make(chan bool)
	}
	stop_chan := make([]chan struct{})
	for i := range stop_chan {
		stop_chan[i] = make(chan bool)
	}
	leave_chan := make([]chan struct{})
	for i := range leave_chan {
		leave_chan[i] = make(chan bool)
	}
	go_chan := make([]chan struct{})
	for i := range go_chan {
		go_chan[i] = make(chan bool)
	}
	Train := func(id int) {
		x_now := time.Now()
		x := time.Since(x_now)
	}
	Gate := func(id int) {}
}
func when(guard bool, channel chan bool) chan bool {
	if !guard {
		return nil
	}
	return channel
}
func time_passage(time_passage []string, ctime time.Duration) int {
	return len(time_passage)
}
