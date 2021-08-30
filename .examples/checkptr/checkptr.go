package main

import "fmt"
import "unsafe"

func main() {
	vals := []int{10, 20, 30, 40}
	start := unsafe.Pointer(&vals[0])
	size := unsafe.Sizeof(int(0))
	for i := 0; i < len(vals); i++ {
		item := *(*int)(unsafe.Pointer(uintptr(start) + size*uintptr(i)))
		fmt.Println(item)
	}
}
