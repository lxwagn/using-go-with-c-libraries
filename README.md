## Using Go with and C (Static and Dynamic) Libraries

#### Lucas Wagner

Since Go has its roots in C and in C's design team at Bell Labs, Go compatibility
with C was important to its designers. Whenever a Go source file imports "C", it 
is using cgo. cgo allows code reuse from C. 

There are countless reasons to reuse C code. Security, simplicity, and saving 
time/money are large ones. Sometimes a developer might want to use a library with Go 
where first-class language support is available for C developers.

Additionally, there are edge cases that inevitably happen with legacy languages like C. 
Sometimes the source code (or the person who wrote the code) cannot be 
found and the compiled library and header file are all that remain. Other 
times, it just makes sense to use tried-and-true code that was created 
and tested over weeks, months, or years.

Go allows linking against C libraries or pasting in in-line code. In this demo, 
both will be shown in the familiar "Hello World" format.

### Quickstart

```
$ make
-------------------------------
Hello from a C library function
Hello from inline C
-------------------------------
$
```

### Step-By-Step

#### Step 1: Build a Shared, Dynamic Library (More Common)

First, compile the source into an object file:

```
gcc -fPIC -c mylib.c
```

Convert the resulting object file(s) into a shared library:

```
gcc -shared -o libmylib.so mylib.o
```

The file src/Makefile contains the full, working code.

##### Building a Static Library (Less Common)

Instead of building a dynamic library, a static library can be created.

First, compile the source into an object file:
```
gcc -c mylib.c
```

Convert the resulting object file(s) into a library:
```
ar rc libmylib.a mylib.o
```

Build an index inside the library:
```
ranlib libmylib.a
```

The file src/Makefile contains the full, working code.

#### Step 2: Add the Go Headers

At this point, you can refer to the provided Go code. Let's take a look at it as a whole 
and then break it down line-by-line:

```
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
```

Let's go through these. Go has some configurable #cgo items, such as
compiler flags (CFLAGS) and linker flags (LDFLAGS). First, we'll need to include 
the `src` directory for the headers we're going to include below:

```#cgo CFLAGS: -I./src```

Next, we're going to include the ./lib directory in our build linkage path and 
in our runtime linkage path. We're going to our shared library, libmylib.so: 

```#cgo LDFLAGS: -L./lib -lmylib -Wl,-rpath=./lib```

We'll need to include some headers so that we have the prototypes for our compiled library:

```#include "mylib.h"```

We are going to be using the C function ```*free()```, and this is contained in stdlib,
so we'll need the stdlib headers.

```#include <stdlib.h>```

Lastly, stdio.h is needed for our in-line function, myPrintFunction2():

```
#include <stdio.h>

void myPrintFunction2() {
	printf("Hello from inline C\n");
}
```

It is important to ensure that the *import "C"* is right after the comment. This is how
Go understands that this is no ordinary comment.

```import "C"```

#### Step 3: Add the Go Code
 
Now comes the best part:
 
```
	// C Library
	mystr := C.CString("Hello from a C library function")
	C.myPrintFunction(mystr)
	C.free(unsafe.Pointer(mystr))

	// Inline C
	C.myPrintFunction2()
```

While the Inline C function is self-explanatory, let's review what's going on with the C Library 
function.

C strings have no character count built into the string. They simply contain a \0 terminator at 
the end of the string, wherever that may be. C has become famous for sloppy developers not
terminating strings with the \0 and their code wandering into unauthorized parts of memory until
it reaches \0.

The designers of Go, coming from a C background, wanted strings to be safe. They built in a 
character count to Go strings so that buffer overflows are not possible.

As such, Go strings are not compatible with C strings (also known as pointers to char (*char)). 
In order to use a C function, you must create a CString. Afterward, you must then free that 
memory once it is no longer needed:
 
```mystr := C.CString("Hello from a C library function")```
 
Then, functions are simply called by their C name:
 
```C.myPrintFunction(mystr)```
 
To complete our cleanup, we use the C standard library function ```free()``` to free the 
memory used by our string. Note the use of unsafe.Pointer(), representing a pointer 
to an arbitrary type:
 
```C.free(unsafe.Pointer(mystr))```
