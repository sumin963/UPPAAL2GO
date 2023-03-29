package main

/*
//global dec
#define N 6					//;
typedef  int id_t;
typedef struct Train{
} Train;
typedef struct Gate{
        id_t list[N+1];		//타입
        int len;
} Gate;
void enqueue(Gate *Gate, id_t element)
{
        Gate->list[Gate->len++] = element;	//Gate->len
}
void dequeue(Gate *Gate )
{
        int i = 0;
        Gate-> len -= 1;
        while (i <Gate-> len)
        {
                Gate->list[i] = Gate->list[i + 1];
                i++;
        }
        Gate->list[i] = 0;
}
id_t front(Gate *Gate )
{
   return Gate->list[0];
}
id_t tail(Gate *Gate )
{
   return Gate->list[Gate->len - 1]; //Gate->len
}
*/
import "C"
import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func main() {
	eps := time.Millisecond * 1
	appr := make([]chan bool, C.N)
	stop := make([]chan bool, C.N)
	leave := make([]chan bool, C.N)
	Go := make([]chan bool, C.N)
	for i := range appr {
		appr[i] = make(chan bool)
		stop[i] = make(chan bool)
		leave[i] = make(chan bool)
		Go[i] = make(chan bool)
	}

	train := func(id C.id_t) {
		//local_val := C.Train{}
		var appr_passage []string
		var Go_passage []string
		var cross_passage []string
		now := time.Now()    //clock t;
		t := time.Since(now) // Cumulative clock t

	safe:
		t = time.Since(now)

		fmt.Println("safe location", id, t)
		appr[id] <- true
		now = time.Now()
		goto appr
	appr:
		t = time.Since(now)

		fmt.Println("appr location", id, t)
		appr_passage = []string{"x==10", "x>10", "x==20", "x>20"}
		switch time_passage(appr_passage, t) {
		case 0:
			goto appr_1
		case 1:
			goto appr_2
		case 2:
			goto appr_3
		case 3:
			goto appr_4
		case 4:
			goto exceptionalLoc
		}
	appr_1:
		t = time.Since(now)

		fmt.Println("appr_1 location", id, t)

		select {
		case <-stop[id]:
			goto stop
		case <-time.After(time.Second*10 - t - eps):
			goto appr_2
		}
	appr_2:
		t = time.Since(now)

		fmt.Println("appr_2 location", id, t)

		select {
		case <-stop[id]:
			goto stop
		case <-time.After(time.Second*10 - t):
			goto appr_3
		case <-time.After(time.Second * 0):
			now = time.Now()

			goto cross
		}
	appr_3:
		t = time.Since(now)

		fmt.Println("appr_3 location", id, t)

		select {
		case <-time.After(time.Second*20 - t - eps):
			goto appr_4
		case <-time.After(time.Second * 0):
			now = time.Now()

			goto cross
		}
	appr_4:
		t = time.Since(now)

		fmt.Println("appr_4 location", id, t)

		select {
		case <-time.After(time.Second*20 - t):
			goto exceptionalLoc
		case <-time.After(time.Second * 0):
			now = time.Now()

			goto cross
		}
	stop:
		t = time.Since(now)

		fmt.Println("stop location", id, t)
		select {
		case Go[id] <- true:
			now = time.Now()
			goto Go
		}
	Go:
		t = time.Since(now)
		Go_passage = []string{"x==7", "x>7", "x==15", "x>15"}
		switch time_passage(Go_passage, t) {
		case 0:
			goto Go_1
		case 1:
			goto Go_2
		case 2:
			goto Go_3
		case 3:
			goto Go_4
		case 4:
			goto exceptionalLoc
		}
	Go_1:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*7 - t - eps):
			goto Go_2
		}
	Go_2:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*7 - t):
			goto Go_3
		case <-time.After(time.Second * 0):
			now = time.Now()

			goto cross
		}
	Go_3:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*15 - t - eps):
			goto Go_4
		case <-time.After(time.Second * 0):
			now = time.Now()

			goto cross

		}
	Go_4:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*7 - t):
			goto exceptionalLoc
		case <-time.After(time.Second * 0):
			now = time.Now()

			goto cross
		}
	cross:
		fmt.Println("cross location", id, t)

		t = time.Since(now)
		fmt.Println("cross location", id, t)
		cross_passage = []string{"x==3", "x>3", "x==5", "x>5"}
		switch time_passage(cross_passage, t) {
		case 0:
			goto cross_1
		case 1:
			goto cross_2
		case 2:
			goto cross_3
		case 3:
			goto cross_4
		case 4:
			goto exceptionalLoc
		}
	cross_1:
		t = time.Since(now)

		fmt.Println("cross_1 location", id, t)

		select {
		case <-time.After(time.Second*3 - t - eps):
			goto cross_2
		}
	cross_2:
		t = time.Since(now)

		fmt.Println("cross_2 location", id, t)

		select {
		case <-time.After(time.Second*3 - t):
			goto cross_3
		case leave[id] <- true:
			goto safe
		}
	cross_3:
		t = time.Since(now)

		fmt.Println("cross_3 location", id, t)

		select {
		case <-time.After(time.Second*5 - t - eps):
			goto cross_4
		case leave[id] <- true:
			goto safe
		}
	cross_4:
		t = time.Since(now)
		fmt.Println("cross_4 location", id, t)

		select {
		case <-time.After(time.Second*5 - t):
			goto exceptionalLoc
		case leave[id] <- true:
			goto safe
		}
	exceptionalLoc:
		fmt.Println("exceptionalLoc", id)
	}

	gate := func() { //selcet부분과, 하나의 로케이션에서 엣지가 여러개일떄 자동으로 생성하는 방법 고려.

		local_val := C.Gate{list: [C.N + 1]C.id_t{0, 0, 0, 0, 0, 0, 0}, len: 0}

	free:
		fmt.Println("gate free")
		select {
		case <-when(local_val.len == 0, appr[0]): //select 수정
			C.enqueue(&local_val, 0)
			goto occ
		case <-when(local_val.len == 0, appr[1]): //select 수정
			C.enqueue(&local_val, 1)
			goto occ
		case <-when(local_val.len == 0, appr[2]): //select 수정
			C.enqueue(&local_val, 2)
			goto occ
		case <-when(local_val.len == 0, appr[3]): //select 수정
			C.enqueue(&local_val, 3)
			goto occ
		case <-when(local_val.len == 0, appr[4]): //select 수정
			C.enqueue(&local_val, 4)
			goto occ
		case <-when(local_val.len == 0, appr[5]): //select 수정
			C.enqueue(&local_val, 5)
			goto occ
		case when(local_val.len > 0, Go[C.front(&local_val)]) <- true:
			goto occ
		}
	occ:
		fmt.Println("gate occ")
		select { //select 전체 수정
		case <-leave[C.front(&local_val)]:
			C.dequeue(&local_val)
			goto free
		case <-appr[0]:
			C.enqueue(&local_val, 0)
			goto annoy
		case <-appr[1]:
			C.enqueue(&local_val, 1)
			goto annoy
		case <-appr[2]:
			C.enqueue(&local_val, 2)
			goto annoy
		case <-appr[3]:
			C.enqueue(&local_val, 3)
			goto annoy
		case <-appr[4]:
			C.enqueue(&local_val, 4)
			goto annoy
		case <-appr[5]:
			C.enqueue(&local_val, 5)
			goto annoy
		}
	annoy:
		select {
		case stop[C.tail(&local_val)] <- true:
			goto occ
		}
	}

	go train(0)
	go train(1)
	go train(2)
	go train(3)
	go train(4)
	go train(5)
	go gate()

	<-time.After(time.Second * 20)
}

// no chan, no condition, 항상 트루인 트랜지션 select 조건, 조건을 만족하면 바운더리내에서 가는게아니라 바로감
// instantaneus loc로 가는 select 조건 수정해야할수도
// uppaal select //

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
	for i, val := range time_passage { // 비교하는거 추가
		if strings.Contains(val, "==") {

			num, _ := strconv.Atoi(val[strings.Index(val, "==")+2:])
			if time.Second*time.Duration(num) > ctime {
				return i
			}
		} else if strings.Contains(val, "<") {
			num, _ := strconv.Atoi(val[strings.Index(val, "<")+1:])

			if time.Second*time.Duration(num) == ctime {
				return i
			}
		}
	}
	return len(time_passage)
}

/*
//global dec
#define  N  6			// const int N = 6;
typedef int id_t;
// local2
typedef struct Local2{           //구조체변환
} Local2;

//local dec
typedef struct Local{           //구조체변환
        id_t list[N+1];
        int len;
} Local;

void enqueue(Local *Local, id_t element)        //구조체 인자로
{
        Local->list[Local->len++] = element;    //구조체 값 사용시 멤버 접근하는 ->사용
}

void dequeue(Local *local)
{
        int i = 0;
      local->len -= 1;
        while (i < local->len)							//&lt; -> < 로 변환
        {
                local->list[i] = local->list[i + 1];
                i++;
        }
        local->list[i] = 0;
}

id_t front(Local *local)
{
   return local->list[0];
}

id_t tail(Local *local)
{
   return local->list[local->len - 1];
}
*/
