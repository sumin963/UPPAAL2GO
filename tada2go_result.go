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
	"time"
)

type chan_t struct{}

func main() {
	appr := make([]chan chan_t, C.N)
	stop := make([]chan chan_t, C.N)
	leave := make([]chan chan_t, C.N)
	Go := make([]chan chan_t, C.N)
	eps := time.Millisecond * 100
	train := func(id int) {
		local_val := C.Train{}
		now := time.Now()    //clock t;
		t := time.Since(now) // Cumulative clock t

	safe:
		fmt.Println("safe location")
		appr[id] <- chan_t{}
		now = time.Now()
		goto appr
	appr:
		fmt.Println("appr location")
		t = time.Since(now)
	appr_1:
		t = time.Since(now)
		select {
		case <-stop[id]:
			goto stop
		case <-time.After(time.Second*10 - t - eps):
			goto appr_2
		}
	appr_2:
		t = time.Since(now)

		select {
		case <-stop[id]:
			goto stop
		case <-time.After(time.Second*10 - t):
			goto appr_3
		case <-time.After(time.Second * 0):
			goto cross
		}
	appr_3:
		t = time.Since(now)

		select {
		case <-stop[id]:
			goto stop
		case <-time.After(time.Second*20 - t - eps):
			goto appr_4
		case <-time.After(time.Second * 0):
			goto cross
		}
	appr_4:
		t = time.Since(now)

		select {
		case <-stop[id]:
			goto stop
		case <-time.After(time.Second*20 - t):
			goto exceptionalLoc
		case <-time.After(time.Second * 0):
			goto cross
		}
	stop:
		t = time.Since(now)
		select {
		case Go[id] <- chan_t{}:
			now = time.Now()
			goto Go
		}
	Go:
		t = time.Since(now)
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
		t = time.Since(now)
	cross_1:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*3 - t - eps):
			goto cross_2
		}
	cross_2:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*3 - t):
			goto cross_3
		case leave[id] <- chan_t{}:
			goto safe
		}
	cross_3:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*5 - t - eps):
			goto cross_4
		case leave[id] <- chan_t{}:
			goto safe
		}
	cross_4:
		t = time.Since(now)
		select {
		case <-time.After(time.Second*5 - t):
			goto exceptionalLoc
		case leave[id] <- chan_t{}:
			goto safe
		}
	exceptionalLoc:
		fmt.Println("exceptionalLoc")
	}

	gate := func() { //selcet부분과, 하나의 로케이션에서 엣지가 여러개일떄 자동으로 생성하는 방법 고려.

		local_val := C.Gate{list: [7]C.id_t{0, 0, 0, 0, 0, 0, 0}, len: 0}

	free:
		select {
		case <-when(local_val.len > 0, appr[0]): //select 수정
			C.enqueue(&local_val, 0)
			goto occ
		case when(local_val.len == 0, Go[C.front(&local_val)]) <- chan_t{}:
			Go[C.front(&local_val)] <- chan_t{}
			goto occ
		}
	occ:
		select { //select 전체 수정
		case <-leave[C.front(&local_val)]:
			C.dequeue(&local_val)
			goto free
		case <-appr[0]:
			C.enqueue(&local_val, 0)
			goto annoy
		}
	annoy:
		select {
		case stop[C.tail(&local_val)] <- chan_t{}:
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

	<-time.After(time.Second * 5)
	fmt.Println()
}

// appr에서 appr_1 appr_2 appr_3 으로 보내는 함수
// no chan, no condition, 항상 트루인 트랜지션 select 조건
// instantaneus loc로 가는 select 조건 수정해야할수도

func when(guard bool, channel chan interface{}) chan interface{} {
	if !guard {
		return nil
	}
	return channel
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
