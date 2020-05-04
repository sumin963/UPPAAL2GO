package main

import (
	"fmt"
	"sync"
	"time"
)

// Global Declarations

// ---------------- Configuration --------------------

const N int = 4 // const int N = 4;          Number of tasks.

type pid_t int // typedef int[0,N-1] pid_t;  Process IDs.
const (
	n0 pid_t = iota
	n1
	n2
	n3
)

var T = [N]int{18, 16, 16, 14} // const int T[pid_t] = { 18, 16, 16, 14 }; End-periods
var E = [N]int{0, 2, 2, 4}     // const int E[pid_t] = {  0,  2,  2,  4 };
var L = [N]int{0, 2, 2, 4}     // const int L[pid_t] = {  0,  2,  2,  4 }; // [ E[i] , L[i] ] Ready interval
var D = [N]int{17, 14, 12, 6}  // const int D[pid_t] = { 17, 14, 12,  6 }; Deadlines
var C = [N]int{6, 2, 4, 5}     // const int C[pid_t] = {  6,  2,  4,  5 }; Computation Times
var P = [N]int{1, 2, 3, 4}     // const int P[pid_t] = {  1,  2,  3,  4 }; Priorities

// Shared resources - semaphores in fact.
const S int = 2 // const int S = 2; Number of semaphores.

type sid_t int // typedef int[0,S-1] sid_t; Resource IDs.
const (
	s0 sid_t = iota
	s1
)

var SP = [S][N]int{ // const int SP[sid_t][pid_t] =  { {1, 0, 0, 2},  {0, 0, 1, 3}};  Resource IDs.
	{1, 0, 0, 2},
	{0, 0, 1, 3},
}

var SV = [S][N]int{ // const int SV[sid_t][pid_t] = { {5, 0, 0, 3}, {0, 0, 3, 4}}; For each task:
	{5, 0, 0, 3},
	{0, 0, 3, 4},
}

// -----------------------------------------------------

func main() {
	var wg sync.WaitGroup
	var priority sync.Mutex
	var committed = 0

	done := make(chan interface{}) //chan	done, ready, run, stop;
	ready := make(chan interface{})
	run := make(chan interface{})
	stop := make(chan interface{})

	var p [N]int       // int p[pid_t]; Dynamic priorities.
	var queue [N]pid_t // pid_t queue[pid_t]; Task queue
	var len int        // int[0,N] len = 0; Length of the queue.

	// Could be const but easier to have it computed.
	var cp [S]int // int cp[sid_t]; Ceiling priorities.

	// Computation intervals for every task.
	// At most (take + release)*S + C for every task.
	// [0] -> bound, [1] -> what to do
	// todo >0 -> take sema todo-1
	// todo <0 -> release sema -todo-1
	// todo ==0 -> end of computation
	// These 2 variables could be const but it's error-prone
	// and not nice to write them manually.
	var ci [N][2*S + 1][2]int //int ci[pid_t][2*S+1][2];
	var ns [N]int             //int ns[pid_t];
	initialize := func() {    //void initialize()
		// Initialize dynamic priorities.
		p = P //p = P;
		// Ceiling priorities.
		for i := 0; i < S; i++ {
			var max int = 0 //int max = 0;
			for j := 0; j < N; j++ {
				if SV[i][j] != 0 { // Task j is using semaphore i.
					max = Max(P[j], max)
				}
			}
			cp[i] = max
		}
		for i := 0; i < N; i++ {
			// Fill.
			var a int // int a,b,elem[2];
			var b int
			var elem [2]int
			ns[i] = 0
			for j := 0; j < S; j++ {
				if SV[j][i] != 0 {
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

			// C[i] always last, no need for ns[i]++ & we count resources.
			// Insertion-sort.

			for a = 1; a < ns[i]; a++ {

				elem = ci[i][a]
				for b = a - 1; b >= 0 && ci[i][b][0] > elem[0]; b-- {

					ci[i][b+1] = ci[i][b]
				}
				ci[i][b+1] = elem
			}
		}
	}
	head := func() pid_t { //pid_t head()
		return queue[0]
	}
	isEmpty := func() bool { //bool isEmpty()
		return len == 0
	}
	remove := func() { //void remove()
		var i int
		for i = 0; i+1 < N; i++ {

			queue[i] = queue[i+1]
		}
		len--
		queue[len] = 0
	}
	// Template Task
	Task := func(id pid_t) {
		// Task Declarations
		axTime := time.Now()     //clock ax, t, wcrt;
		ax := time.Since(axTime) //Cumulative clock ax
		tTime := time.Now()
		t := time.Since(tTime) //Cumulative clock t
		//wcrtTime:=time.Now()

		var r int = 0    //int r = 0;
		var sema [S]bool //bool sema[sid_t];

		add := func() { //void add()
			var i int //pid_t i, tmp;
			var tmp pid_t
			queue[len] = id

			for i = len; i > 0 && p[queue[i]] > p[queue[i-1]]; i-- {

				tmp = queue[i]
				queue[i] = queue[i-1]
				queue[i-1] = tmp
			}
			len++
		}

		updatePriority := func(s int) { //void updatePriority(int s)
			if s > 0 { // Take.
				s = s - 1
				sema[s] = true
				p[id] = Max(cp[s], p[id])
			} else { // Release
				var j int
				var tmp pid_t

				s = -s - 1
				sema[s] = false
				// Recompute priority.

				p[id] = P[id]
				for i := 0; i < S; i++ {
					if sema[i] {
						p[id] = Max(cp[i], p[id])
					}
				}
				// Reorder.
				for j = 0; j+1 < len && (p[queue[j]] < p[queue[j+1]]); j++ {

					tmp = queue[j]
					queue[j] = queue[j+1]
					queue[j+1] = tmp
				}
			}
		}

	Idle: // Idle Location
		for {
			fmt.Println(id, "Task  idle")
			time.Sleep(time.Second * time.Duration(E[id]))
			//invariant

			ready <- struct{}{} // [Idle --> Ready] | Sync: ready!
			tTime = time.Now()  // [Idle --> Ready] | Update: t=0
			//wcrtTime = time.Now()          // [Idle --> Ready] | Update: wcrt=0

			add() // [Idle --> Ready] | Update: add()
			goto Ready
		}
	Ready: // Ready Location
		fmt.Println(id, "Task  Ready")
		for {
			t = time.Since(tTime) // [Ready --> Running] | Update: t=0
			switch {
			case head() == id: // [Ready --> Running] | Guard: head()==0
				<-run               // [Ready --> Running] | Sync: run?
				axTime = time.Now() // [Ready --> Running] | Update: ax=0
				goto Running
			case (int(t) / 1000000000) > D[id]: // [Ready --> Error] | Guard: t>D[id]
				committed++ // Committed Location
				priority.Lock()
				goto Error
			default:
			}
		}

	Running: // Running Location
		fmt.Println(id, "Task  Running")
		for {
			ax = time.Since(axTime)
			t = time.Since(tTime)
			switch {
			case r < ns[id] && (int(ax)/1000000000) == ci[id][r][0]: // [Running --> Running] | Guard: r < ns[id] && ax == ci[id][r][0]
				select {
				case <-stop: // [Running --> Blocked] | sync: stop?
					goto Blocked
				case ready <- struct{}{}: // [Running --> Running] | Sync: ready!
					updatePriority(ci[id][r][1]) // [Running --> Running] | Update: updatePriority()
					r++                          // [Running --> Running] | Update: r++
					goto Running
				default:
				}
			case (head() == id && (int(ax)/1000000000) >= C[id]) && r == ns[id]: // [Running --> EndPeriod] | Guard: head() == id && ax>=C[id] && r == ns[id]
				done <- struct{}{} // [Running --> EndPeriod] | Sync: done!
				remove()           // [Running --> EndPeriod] | Update: remove()
				r = 0              // [Running --> EndPeriod] | Update: r=0
				goto EndPeriod
			case (int(t) / 1000000000) > D[id]: // [Running --> Error] | Guard: t>D[id]
				committed++ // Committed Location
				priority.Lock()
				goto Error
			case (int(ax) / 1000000000) > ci[id][r][0]: // [Running] | invariant: ax<=ci[id][r][0]
				committed++ // Committed Location
				priority.Lock()
				goto Error
			default:
				select {
				case <-stop: // [Running --> Blocked] | sync: stop?
					goto Blocked
				default:
				}
			}
		}
	Blocked: // Blocked Location
		fmt.Println(id, "Task  Blocked")
		for {
			t = time.Since(tTime)
			switch {
			case head() == id: // [Blocked --> Running] | Guard: head()==id
				<-run // [Blocked --> Running] | sync: run?
				ax_2 := time.Since(axTime)
				axTime = axTime.Add(ax_2 - ax)
				goto Running
			case (int(t) / 1000000000) > D[id]: // [Blocked --> Error] | Guard: t>D[id]
				committed++ // Committed Location
				priority.Lock()
				goto Error
			default:
			}
		}
	EndPeriod: // EndPeriod Location
		fmt.Println(id, "Task  EndPeriod")
		for {
			t = time.Since(tTime)
			switch {
			case (int(t) / 1000000000) == T[id]: // [EndPeriod --> Idle] | Guard: t==T[id]
				tTime = time.Now() // [EndPeriod --> Idle] | Update: t=0
				//wcrt                              // [EndPeriod --> Idle] | Update: wcrt=0
				//&& wcrt'==0:
				goto Idle
			case (int(t) / 1000000000) > T[id]: // [EndPeriod --> Idle] | Invariant: t<=T[id]
				committed++ // Committed Location
				priority.Lock()
				goto Error
			}
		}

	Error: // Error Location // <Alarm>
		fmt.Println(id, "Task  Error")
		wg.Done()
		committed--
		priority.Unlock()
	}

	// Template Task
	Scheduler := func() {
		committed++ // Committed Location
		priority.Lock()
		fmt.Println("Scheduler Init")
		initialize() // [Init --> Free] | Update: initialize()
		committed--
		priority.Unlock()
	Free: // Free Location
		committed++ // Committed Location
		priority.Lock()
		fmt.Println("Scheduler Free")
		switch {
		case isEmpty() == true: // [Free --> anonymous] | Guard: isEmpty()
			committed--
			priority.Unlock()
			fmt.Println("Scheduler anonymous")
			<-ready // [anonymous --> Select] | Sync: ready?
		case isEmpty() == false:
			committed--
			priority.Unlock()
		}
	Select: // Select Location
		committed++ // Committed Location
		priority.Lock()
		fmt.Println("Scheduler Select")
		run <- struct{}{} // [Select --> Occ] | Sync: run!
		committed--
		priority.Unlock()
		goto Occ
	Occ: // Occ Location
		fmt.Println("Scheduler Occ")
		select {
		case <-done: // [Occ --> Free] | Sync: done?
			fmt.Println("Scheduler done")
			goto Free
		case <-ready: // [Occ --> anonnymous] | Sync: ready?
			committed++ // Committed Location
			priority.Lock()
			fmt.Println("Scheduler ready")
			stop <- struct{}{} // [anonnymous --> Select] | Sync: stop?
			committed--
			priority.Unlock()
			goto Select
		}
	}
	wg.Add(4)
	// System declarations

	go Scheduler() //Scheduler
	go Task(n0)    //Task(0)
	go Task(n1)    //Task(1)
	go Task(n2)    //Task(2)
	go Task(n3)    //Task(3)

	wg.Wait()
}

func Max(x int, y int) int { // operator >?
	if x < y {
		return y
	}
	return x
}
