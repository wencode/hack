package cgocall

//#include <stdlib.h>
import "C"

func Abs(i int) int {
	n, _ := C.abs(C.int(i))
	return int(n)
}
