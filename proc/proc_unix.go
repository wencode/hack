// +build linux darwin

package proc

import (
	"unsafe"
)

func Call(fn uintptr, args ...uintptr) uintptr {
	ta := new(tramargs)
	ta.fn = fn
	copy(ta.arg[:], args)
	runtime_cgocall(unsafe.Pointer(trampoline_fn), uintptr(unsafe.Pointer(ta)))
	return ta.ret
}

func SmallCall(smallfn uintptr, args ...uintptr) uintptr {
	ta := new(tramargs)
	ta.fn = smallfn
	copy(ta.arg[:], args)
	realcall(ta)
	return ta.ret
}

type tramargs struct {
	fn uintptr
	arg [6]uintptr
	ret  uintptr
	pad [4]byte
}


var (
	trampoline_fn uintptr
)

func init() {
	trampoline_fn = trampoline_addr()
}

func trampoline(*tramargs)
func trampoline_addr() uintptr
func realcall(*tramargs)

//go:linkname runtime_cgocall runtime.cgocall
func runtime_cgocall(unsafe.Pointer, uintptr) int32
