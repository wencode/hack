package proc

import (
	"syscall"
)

func Call(fn uintptr, a ...uintptr) uintptr {
	var ret uintptr = 0xFFFFFFFFFFFFFFFF
	switch len(a) {
	case 0:
		ret, _, _ = syscall.Syscall(fn, uintptr(len(a)), 0, 0, 0)
	case 1:
		ret, _, _ = syscall.Syscall(fn, uintptr(len(a)), a[0], 0, 0)
	case 2:
		ret, _, _ = syscall.Syscall(fn, uintptr(len(a)), a[0], a[1], 0)
	case 3:
		ret, _, _ = syscall.Syscall(fn, uintptr(len(a)), a[0], a[1], a[2])
	case 4:
		ret, _, _ = syscall.Syscall6(fn, uintptr(len(a)), a[0], a[1], a[2], a[3], 0, 0)
	case 5:
		ret, _, _ = syscall.Syscall6(fn, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], 0)
	case 6:
		ret, _, _ = syscall.Syscall6(fn, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5])
	default:
		panic("call() with too many arguments")
	}
	return ret
}
