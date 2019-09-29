package main

import "C"
import (
	"runtime"
	"log"
	"unsafe"

	"github.com/wencode/hack/dl"
	"github.com/wencode/hack/proc"
)

var (
	libcfile string
)

func init() {
	switch runtime.GOOS {
	case "darwin":
		libcfile = "/usr/lib/libc.dylib"
	case "windows":
	default:
		libcfile = "libc.so.6"
	}

}


func main() {
	lib, err := dl.Open(libcfile)
	if err != nil {
		log.Fatal(err)
	}
	defer lib.Close()
	write_addr := lib.Sym("write")
	if write_addr == 0 {
		log.Fatalf("can't find printf symbol")
	}

	str := []byte("hello,world\n")
	l := proc.Call(write_addr, 1, uintptr(unsafe.Pointer(&(str[0]))), uintptr(len(str)))
	if l != len(str) {
		log.Fatal("call write failed")
	}
}
