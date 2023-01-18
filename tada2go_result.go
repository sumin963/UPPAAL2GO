package main

/*
#define N 6			// const int N = 6;
typedef int id_t;


id_t list[N+1];
int len;

void enqueue(id_t element)
{
        list[len++] = element;
}

void dequeue()
{
        int i = 0;
        len -= 1;
        while (i < len)
        {
                list[i] = list[i + 1];
                i++;
        }
        list[i] = 0;
}

id_t front()
{
   return list[0];
}

id_t tail()
{
   return list[len - 1];
}
*/
import "C"
import "fmt"

func main() {
	fmt.Println(C.list)
}

/*
const int N = 6;
typedef int[0,N-1] id_t;					//c스타일로 변형**

chan        appr[N], stop[N], leave[N];		//golang 글로벌로 빼서 선언
urgent chan go[N];							//golang 글로벌로 빼서 선언

//local
clock x;									//golang 글로벌로 빼서 선언


id_t list[N+1];								//변형안해도됨 typedef만 제대로되면
int[0,N] len;								//c스타일로 변형**

void enqueue(id_t element)
{
        list[len++] = element;
}

void dequeue()
{
        int i = 0;
        len -= 1;
        while (i &lt; len)
        {
                list[i] = list[i + 1];
                i++;
        }
        list[i] = 0;
}

id_t front()
{
   return list[0];
}

id_t tail()
{
   return list[len - 1];
}
*/
