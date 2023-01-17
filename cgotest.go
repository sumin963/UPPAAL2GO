package main

/*
int list[5]={1,2,3,4,5};
int len=3;

int tail()
{
   return list[len - 1];
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

void enqueue(int element)
{
        list[len++] = element;
}

*/
import "C"
import "fmt"

func main() {

	fmt.Println(C.enqueue(3), C.list)
	// Output: 42
}

/*
const int N = 6;         // # trains
typedef int[0,N-1] id_t;

id_t list[N+1];
int[0,N] len;

// Put an element at the end of the queue
void enqueue(id_t element)
{
        list[len++] = element;
}

// Remove the front element of the queue
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

// Returns the front element of the queue
id_t front()
{
   return list[0];
}

// Returns the last element of the queue
id_t tail()
{
   return list[len - 1];
}
*/

/*
const int N = 6;         // # trains
typedef int[0,N-1] id_t;

chan        appr[N], stop[N], leave[N];
urgent chan go[N];


clock x;


id_t list[N+1];
int[0,N] len;

// Put an element at the end of the queue
void enqueue(id_t element)
{
        list[len++] = element;
}

// Remove the front element of the queue
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

// Returns the front element of the queue
id_t front()
{
   return list[0];
}

// Returns the last element of the queue
id_t tail()
{
   return list[len - 1];
}
*/
