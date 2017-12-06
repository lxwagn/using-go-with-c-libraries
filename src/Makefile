# Places a static or shared (dynamic) library into ../lib

all: dynamic
	
dynamic:
	gcc -fPIC -c mylib.c
	gcc -shared -o libmylib.so mylib.o
	mv -f libmylib.so ../lib
	rm -f mylib.o

static:
	gcc -c mylib.c
	ar rc libmylib.a mylib.o
	ranlib libmylib.a
	mv -f libmylib.a ../lib
	rm -f mylib.o

