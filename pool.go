package main

/*
#include <stdio.h>
#include "stdlib.h"

void getStr(char *str) {
    sprintf(str, "hello world");
}
*/
import "C"
import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"
)

var (
	ObjectPool  *sync.Pool
	alloc, free uint64
)

type PoolObject struct {
	data *C.char // *C.char
}

func initObjectPool(len int) {
	ObjectPool = &sync.Pool{
		New: func() interface{} {
			b := make([]byte, len)
			obj := &PoolObject{
				data: (*C.char)(C.CBytes(b)),
			}
			atomic.AddUint64(&alloc, 1)
			runtime.SetFinalizer(obj, func(object *PoolObject) {
				if object.data != nil {
					C.free((unsafe.Pointer)(object.data))
					object.data = nil

				}
				atomic.AddUint64(&free, 1)
			})
			return obj
		},
	}
}

func main() {
	poolObjectDataLen := 1024
	initObjectPool(poolObjectDataLen)

	str := ObjectPool.Get().(*PoolObject)

	C.getStr(str.data)
	// print data
	fmt.Println(string((*[1 << 30]byte)(unsafe.Pointer(str.data))[:poolObjectDataLen]))

	ObjectPool.Put(str)
}
