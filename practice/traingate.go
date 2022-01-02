package main

import (
	"fmt"
	"math/rand"
	"time"
)

type id_t int
type chan_t struct{}

const (
	n0 id_t = iota
	n1
	n2
	n3
	n4
	n5
)
const N int = 6

func main() {
	appr := make([]chan chan_t, N)
	stop := make([]chan chan_t, N)
	leave := make([]chan chan_t, N)
	Go := make([]chan chan_t, N)
	id := []id_t{n0, n1, n2, n3, n4, n5}
	for _, x := range id {
		go train(appr, stop, leave, Go, x)
	}
	go gate(appr, leave, stop, Go)

	time.Sleep(time.Second * 10000)
	fmt.Println("main Done")
}
func train(appr, leave, stop, Go []chan chan_t, id id_t) {
	for {
		fmt.Printf("train %d Safe\n", id)
		appr[id] <- chan_t{}
		fmt.Printf("train %d Sae\n", id)
		select {
		case <-stop[id]:
			fmt.Printf("train %d Stop\n", id)
			<-Go[id]
			fmt.Printf("train %d Start\n", id)
			time.Sleep(time.Second * 7)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(8000)))
			fmt.Printf("train %d Cross", id)
			time.Sleep(time.Second * 3)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
			leave[id] <- chan_t{}
			continue

		case <-time.After(time.Second * 10):
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(10000)))
			fmt.Printf("train %d Cross", id)
			time.Sleep(time.Second * 3)
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(2000)))
			leave[id] <- chan_t{}
			continue
		}
	}
}
func gate(appr, leave, stop, Go []chan chan_t) {
	var list [N + 1]id_t
	len := 0
	dequeue := func() {
		var i int
		len--
		for {
			if i < len {
				break
			}
			list[i] = list[i+1]
			i++
		}
		list[i] = 0
	}
	enqueue := func(element id_t) {
		len++
		list[len] = element

	}
	front := func() id_t {
		return list[0]
	}
	tail := func() id_t {
		return list[len-1]
	}

	for {
		if len == 0 {
			select {
			case <-appr[n0]:
				enqueue(n0)
				for {
					select {
					case <-appr[n0]:
						enqueue(n0)
						stop[tail()] <- chan_t{}
					case <-appr[n1]:
						enqueue(n1)
						stop[tail()] <- chan_t{}
					case <-appr[n2]:
						enqueue(n2)
						stop[tail()] <- chan_t{}
					case <-appr[n3]:
						enqueue(n3)
						stop[tail()] <- chan_t{}
					case <-appr[n4]:
						enqueue(n4)
						stop[tail()] <- chan_t{}
					case <-appr[n5]:
						enqueue(n5)
						stop[tail()] <- chan_t{}
					case <-leave[front()]:
						dequeue()
						break
					}
				}
			}
		} else if len > 0 {
			select {
			case <-appr[front()]:
				Go[front()] <- chan_t{}
				for {
					select {
					case <-appr[n0]:
						enqueue(n0)
						stop[tail()] <- chan_t{}
					case <-appr[n1]:
						enqueue(n1)
						stop[tail()] <- chan_t{}
					case <-appr[n2]:
						enqueue(n2)
						stop[tail()] <- chan_t{}
					case <-appr[n3]:
						enqueue(n3)
						stop[tail()] <- chan_t{}
					case <-appr[n4]:
						enqueue(n4)
						stop[tail()] <- chan_t{}
					case <-appr[n5]:
						enqueue(n5)
						stop[tail()] <- chan_t{}
					case <-leave[front()]:
						dequeue()
						break
					}
				}
			}
		}
	}

}
