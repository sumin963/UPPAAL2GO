package main

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
import "C"
import (
	"fmt"
	"time"
)

func main() {
	now := time.Now()
	routine := func() {
		point := C.Local{list: [7]C.id_t{0, 0, 0, 0, 0, 0, 0}, len: 3}

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
