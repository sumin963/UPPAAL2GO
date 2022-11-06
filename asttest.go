package main

import "fmt"

func main() {
	src := `
/*
 * For more details about this example, see 
 * "Automatic Verification of Real-Time Communicating Systems by Constraint Solving", 
 * by Wang Yi, Paul Pettersson and Mats Daniels. In Proceedings of the 7th International
 * Conference on Formal Description Techniques, pages 223-238, North-Holland. 1994.
 */

const int N = 6;         // # trains
typedef int[0,N-1] id_t;

chan        appr[N], stop[N], leave[N];
urgent chan go[N];
"
`
	fmt.Println(src)

}
