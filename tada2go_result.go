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

	train := func() {
		local_val := C.Train{}
		now := time.Now()    //clock t;
		t := time.Since(now) // Cumulative clock t

	}
	gate := func() {
		local_val := C.Gate{list: [7]C.id_t{0, 0, 0, 0, 0, 0, 0}, len: 0}

	}
	now := time.Now()
	routine := func() {
		point := C.Gate{list: [7]C.id_t{0, 0, 0, 0, 0, 0, 0}, len: 3}

		C.enqueue(&point, 1)
		C.enqueue(&point, 1)
		fmt.Println(C.front(&point))
		fmt.Println(C.tail(&point))
		fmt.Println(point.list)
		C.dequeue(&point)
		fmt.Println(point.list, point.len)
	}
	go routine()
	go routine()

	<-time.After(time.Second * 5)
	t := time.Since(now)
	fmt.Println(t)
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
