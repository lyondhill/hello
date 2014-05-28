package main

/*
#include <stdio.h>
#include <stdlib.h>

void myprint(char* s) {
        printf("(cpart): %s", s);
}
*/
import "C"

import "unsafe"

func Example() {
	cs := C.CString("(gopart)Hello from go\n")
	C.myprint(cs)
	C.free(unsafe.Pointer(cs))
}

func main() {
	Example()
}
