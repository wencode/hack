// +build linux darwin

package proc

import (
	"unsafe"
)

func Call(fn uintptr, args ...uintptr) uintptr {
	ta := new(tramargs)
	ta.fn = fn
	copy(ta.arg[:], args)
	call(ta)
	return ta.ret
}

func call(a *tramargs) int32 {
	return runtime_cgocall(unsafe.Pointer(trampoline_fn), uintptr(unsafe.Pointer(a)))
}

type tramargs struct {
	fn uintptr
	arg [6]uintptr
	ret  uintptr
}


var (
	trampoline_fn uintptr
)

func init() {
	trampoline_fn = trampoline_addr()
}

func trampoline(*tramargs)
func trampoline_addr() uintptr

//go:linkname runtime_cgocall runtime.cgocall
func runtime_cgocall(unsafe.Pointer, uintptr) int32
