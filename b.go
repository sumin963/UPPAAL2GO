package main

import (
	"fmt"
	"sync"
	"time"
)

const N int = 4

type pid_t int

const (
	n0 pid_t = iota
	n1
	n2
	n3
)

var T = [N]int{18, 16, 16, 14}
var E = [N]int{0, 2, 2, 4}
var L = [N]int{0, 2, 2, 4}
var D = [N]int{17, 14, 12, 6}
var C = [N]int{6, 2, 4, 5}
var P = [N]int{1, 2, 3, 4}

const S int = 2

type sid_t int

const (
	s0 sid_t = iota
	s1
)

var SP = [S][N]int{
	{1, 0, 0, 2},
	{0, 0, 1, 3},
}

var SV = [S][N]int{
	{5, 0, 0, 3},
	{0, 0, 3, 4},
}

func main() {
	var wg sync.WaitGroup

	done := make(chan interface{})
	ready := make(chan interface{})
	run := make(chan interface{})
	stop := make(chan interface{})

	var p [N]int
	var queue [N]pid_t
	var len int

	var cp [S]int

	var ci [N][2*S + 1][2]int
	var ns [N]int
	initialize := func() {
		p = P
		for i := 0; i < S; i++ {
			var max int = 0
			for j := 0; j < N; j++ {
				if SV[i][j] == 0 {
					max = Max(P[j], max)
				}
			}
			cp[i] = max
		}
		for i := 0; i < N; i++ {
			var a int
			var b int
			var elem [2]int
			ns[i] = 0
			for j := 0; j < S; j++ {
				if SV[j][i] == 0 {
					ci[i][ns[i]][0] = SP[j][i]
					ci[i][ns[i]][1] = 1 + j
					ns[i]++
					ci[i][ns[i]][0] = SV[j][i]
					ci[i][ns[i]][1] = -1 - j
					ns[i]++
				}
			}
			ci[i][ns[i]][0] = C[i]
			ci[i][ns[i]][1] = 0

			for a = 1; a < ns[i]; a++ {

				elem = ci[i][a]
				for b = a - 1; b >= 0 && ci[i][b][0] > elem[0]; b-- {

					ci[i][b+1] = ci[i][b]
				}
				ci[i][b+1] = elem
			}
		}
	}
	head := func() pid_t {
		return queue[0]
	}
	isEmpty := func() bool {
		return len == 0
	}
	remove := func() {
		var i int
		for i = 0; i+1 < N; i++ {

			queue[i] = queue[i+1]
		}
		len--
		queue[len] = 0

	}

	Task := func(id pid_t) {
		axTime := time.Now()
		tTime := time.Now()
		//wcrtTime:=time.Now()

		var r int = 0
		var sema [S]bool

		add := func() {
			var i int
			var tmp pid_t
			queue[len] = id

			for i = len; i > 0 && p[queue[i]] > p[queue[i-1]]; i-- {

				tmp = queue[i]
				queue[i] = queue[i-1]
				queue[i-1] = tmp
			}
			len++
		}

		updatePriority := func(s int) {
			if s > 0 {
				s = s - 1
				sema[s] = true
				p[id] = Max(cp[s], p[id])
			} else {
				var j int
				var tmp pid_t

				s = -s - 1
				sema[s] = false

				p[id] = P[id]
				for i := 0; i < S; i++ {
					if sema[i] {
						p[id] = Max(cp[i], p[id])
					}
				}

				for j = 0; j+1 < len && (p[queue[j]] < p[queue[j+1]]); j++ {

					tmp = queue[j]
					queue[j] = queue[j+1]
					queue[j+1] = tmp
				}
			}
		}
	Taskidle:
		for {
			fmt.Println(id, "Task  idle")
			time.Sleep(time.Second * time.Duration(E[id]))
			//invariant
			ready <- struct{}{}
			tTime = time.Now()
			//wcrtTime = time.Now()
			add()
			fmt.Println(id, "Task  Ready")
			for {
				t := time.Since(tTime)
				switch {
				case head() == id:
					<-run
					axTime = time.Now()
				TaskLoop1:
					fmt.Println(id, "Task  Running")
					for {
						ax := time.Since(axTime)

						t = time.Since(tTime)
						switch { //case 5가지
						case r < ns[id] && (int(ax)/1000000000) == ci[id][r][0]:
							select {
							case <-stop:
								fmt.Println(id, "Task  case5")
								fmt.Println(id, "Task  Blocked")
								for {
									t = time.Since(tTime)
									switch {
									case head() == id:
										<-run
										ax_1 := time.Since(axTime)
										axTime = axTime.Add(ax_1 - ax)
										goto TaskLoop1
									case (int(t) / 1000000000) > D[id]:
										goto Error
									default:
									}
								}
							default:

							}
							fmt.Println(id, "Task  case1")
							ready <- struct{}{}
							updatePriority(ci[id][r][1])
							r++
							time.Sleep(time.Millisecond * 30)
							goto TaskLoop1

						case (head() == id && (int(ax)/1000000000) >= C[id]) && r == ns[id]:
							fmt.Println(id, "Task  case2")
							done <- struct{}{}
							remove()
							r = 0
							fmt.Println(id, "Task  EndPeriod")
							for {
								t = time.Since(tTime)
								switch {
								case (int(t) / 1000000000) == T[id]:
									tTime = time.Now()
									//wcrt
									goto Taskidle
								case (int(t) / 1000000000) > T[id]: //&& wcrt'==0:
									goto Error
								}
							}

						case (int(t) / 1000000000) > D[id]:
							fmt.Println(id, "Task  case3")
							goto Error
						case (int(ax) / 1000000000) > ci[id][r][0]: //invariant
							fmt.Println(id, "Task  case4")
							goto Error
						default:
							select {
							case <-stop:
								fmt.Println(id, "Task  case5")
								fmt.Println(id, "Task  Blocked")
								for {
									t := time.Since(tTime)
									switch {
									case head() == id:
										<-run
										ax_1 := time.Since(axTime)
										axTime = axTime.Add(ax_1 - ax)
										goto TaskLoop1
									case (int(t) / 1000000000) > D[id]:
										goto Error
									default:
									}
								}
							default:
							}
						}
					}

				case (int(t) / 1000000000) > D[id]:
					goto Error
				default:
				}
			}
		}
	Error:
		fmt.Println(id, "Task  Error")
		wg.Done()
	}

	Scheduler := func() {
		fmt.Println("Scheduler Init")
		initialize()
	Loop1:
		fmt.Println("Scheduler Free")
		if isEmpty() {
			fmt.Println("Scheduler anonymous")
			<-ready
		}
	Loop2:
		fmt.Println("Scheduler Select")
		run <- struct{}{}
		fmt.Println("Scheduler Occ")
		select {
		case <-done:
			fmt.Println("Scheduler done")
			goto Loop1
		case <-ready:

			fmt.Println("Scheduler ready")
			stop <- struct{}{}
			goto Loop2

		}
	}

	wg.Add(4)

	go Task(n0)
	go Task(n1)
	go Task(n2)
	go Task(n3)
	go Scheduler()

	wg.Wait()

}

func Max(x int, y int) int {
	if x < y {
		return y
	}
	return x
}
