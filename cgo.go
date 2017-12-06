package main

/*

#cgo CFLAGS: -I./src
#cgo LDFLAGS: -L./lib -lmylib -Wl,-rpath=./lib
#include "mylib.h"
#include <stdlib.h>
#include <stdio.h>

void myPrintFunction2() {
	printf("Hello from inline C\n");
}

*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {

	fmt.Println("-------------------------------")

	// C Library
	mystr := C.CString("Hello from a C library function")
	C.myPrintFunction(mystr)
	C.free(unsafe.Pointer(mystr))

	// Inline C
	C.myPrintFunction2()

	fmt.Println("-------------------------------")
}
