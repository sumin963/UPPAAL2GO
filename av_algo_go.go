package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// #define pubsubtime1
// #define pub_num4
// #define sub_num1
// #define sen_num4
// typedef  int pub;
// typedef  int sen_pub;
// typedef  int sub;
// typedef  int pub_len;
// #define ctimemincons
// #define ctimemaxcons
// #define periodcons
// #define N 10
// typedef  int id_t;
// typedef struct Node{
// }Node;
// typedef struct Controller{
// }Controller;
//
import "C"

func main() {
	eps := time.Millisecond * 10
	queue_chan := make([]chan bool, C.pub)
	for i := range queue_chan {
		queue_chan[i] = make(chan bool)
	}
	sen_queue_chan := make([]chan bool, C.sen_pub)
	for i := range sen_queue_chan {
		sen_queue_chan[i] = make(chan bool)
	}
	Node := func(pub_id C.pub) {
		x_now := time.Now()
		x := time.Since(x_now)
		t_now := time.Now()
		t := time.Since(t_now)
		// var id0_passage []string
		// var id1_passage []string
		// var id2_passage []string
		// var id3_passage []string
		// var id5_passage []string
		// var id6_passage []string
		// var id7_passage []string
		// var id8_passage []string
	id0:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id0", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*0 - C.sen_len[pub_id] - eps):
			goto id3
		case <-time.After(time.Second*0 - C.sen_len[pub_id]):
			goto id1
		}
	id1:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id1", "location", "x", ":", x)
		select {
		case <-sen_queue_chan[pub_id]:
			goto id1
		case <-time.After(time.Second*0 - x - eps):
			goto id2
		}
	id2:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id2", "location", "x", ":", x)
		select {
		case <-sen_queue_chan[pub_id]:
			goto id2
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			goto id3
		case <-time.After(time.Second*0 - x):
			goto id4
		}
	id3:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id3", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*0 - x - eps):
			goto id5
		}
	id4:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "exp", "location", "x", ":", x)
		select {}
	id5:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id5", "location", "x", ":", x)
		select {
		case <-time.After(time.Second * 0):
			goto id9
		case <-time.After(time.Second*0 - x):
			goto id6
		}
	id6:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id6", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*0 - x):
			goto id4
		case <-time.After(time.Second * 0):
			goto id9
		}
	id7:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id7", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*0 - t - eps):
			goto id8
		}
	id8:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id8", "location", "x", ":", x)
		select {
		case <-time.After(time.Second*0 - t):
			goto id4
		case <-time.After(time.Second * 0):
			t_now = time.Now()
			goto id0
		}
	id9:
		x = time.Since(x_now)
		fmt.Println("Node", "template", "id9", "location", "x", ":", x)
		select {
		case queue_chan[pub_id] <- true:
			goto id7
		}
	}
	Controller := func() {
		x_now := time.Now()
		x := time.Since(x_now)
		t_now := time.Now()
		t := time.Since(t_now)
	id10:
		select {
		case <-time.After(time.Second*0 - t):
			goto id11
		case <-time.After(time.Second * 0):
			t_now = time.Now()
			goto id13
		}
	id11:
		select {}
	id12:
		select {
		case <-time.After(time.Second*0 - x):
			goto id11
		case <-time.After(time.Second * 0):
			goto id14
		}
	id13:
		select {
		case <-time.After(time.Second*0 - C.sen_len[3] - eps):
			goto id17
		case <-time.After(time.Second*0 - C.sen_len[3]):
			goto id19
		}
	id14:
		select {
		case sen_queue_chan[3] <- true:
			goto id16
		}
	id15:
		select {
		case <-time.After(time.Second * 0):
			goto id14
		case <-time.After(time.Second*0 - x):
			goto id12
		}
	id16:
		select {
		case <-time.After(time.Second*0 - t - eps):
			goto id10
		}
	id17:
		select {
		case <-time.After(time.Second*0 - x - eps):
			goto id15
		}
	id18:
		select {
		case <-queue_chan[0]:
			goto id18
		case <-queue_chan[2]:
			goto id18
		case <-queue_chan[1]:
			goto id18
		case <-time.After(time.Second * 0):
			x_now = time.Now()
			goto id17
		case <-time.After(time.Second*0 - x):
			goto id11
		}
	id19:
		select {
		case <-queue_chan[2]:
			goto id19
		case <-queue_chan[0]:
			goto id19
		case <-queue_chan[1]:
			goto id19
		case <-time.After(time.Second*0 - x - eps):
			goto id18
		}
	}
	go Node(0)
	go Node(1)
	go Node(2)
	go Controller()
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
