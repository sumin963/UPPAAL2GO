package main

import "fmt"

func main() {

	num := []int{4, 5, 6}
	for i, v := range num {
		fmt.Println(i, v)
	}
}
