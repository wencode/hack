package proc

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/wencode/hack/dl"
	"github.com/wencode/hack/proc"

	"github.com/wencode/hack/test/cgocall"
)

var (
	libcfile     string
	write_symbol string
)

func init() {
	switch runtime.GOOS {
	case "darwin":
		libcfile = "/usr/lib/libc.dylib"
		write_symbol = "write"
	case "windows":
		libcfile = "msvcrt.dll"
		write_symbol = "_write"
	default:
		libcfile = "libc.so.6"
		write_symbol = "write"
	}

}

func TestWrite(t *testing.T) {
	lib, err := dl.Open(libcfile)
	if err != nil {
		t.Fatal(err)
	}
	defer lib.Close()
	write_addr := lib.Sym(write_symbol)
	if write_addr == 0 {
		t.Fatalf("can't find write symbol")
	}

	str := []byte("hello,world\n")
	l := proc.Call(write_addr, 1, uintptr(unsafe.Pointer(&(str[0]))), uintptr(len(str)))
	if int(l) != len(str) {
		t.Errorf("call write failed")
	}
}

func TestExit(t *testing.T) {
	lib, err := dl.Open(libcfile)
	if err != nil {
		t.Fatal(err)
	}
	defer lib.Close()
	exit_addr := lib.Sym("exit")
	if exit_addr == 0 {
		t.Fatalf("can't find exit symbol")
	}
	proc.Call(exit_addr, 0)
}

func BenchmarkAbsCGO(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cgocall.Abs(i)
	}
}

func BenchmarkAbsProc(b *testing.B) {
	b.StopTimer()
	lib, err := dl.Open(libcfile)
	if err != nil {
		b.Fatal(err)
	}
	defer lib.Close()
	abs_addr := lib.Sym("abs")
	if abs_addr == 0 {
		b.Fatal("can't find abs symbol")
	}

	b.StartTimer()
	for i := 0; i < b.N; i++ {
		proc.Call(abs_addr, uintptr(i))
	}
}
