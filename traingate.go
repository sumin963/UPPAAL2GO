package main

import (
	"sync"
	"time"
)
type id_t int

const (
	n0 =1
	n1 =2
	n2 =3
	n3 =4
	n4 =5
	n5 =6
)
const(
	N int=6
)



var appr chan []int
var leave chan []int
var stop chan []int
//var go chan []int

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
}
func train(){
	appr<-interface{}
	select {
	case <-stop:
	case<-time.After(time.Second*10):
	}
}
func gate(){
	var list [N+1]int
	len :=len(list)
}
func dequeue(list []int,len int){
	var i int
	len--
	for{
		if i<len{
			break
		}
		list[i] = list[i+1]
		i++
	}
	list[i]=0
}
func enqueue(list []int,len int,element id_t){
	list[len++]=element
}