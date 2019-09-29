package proc

func Call(fn uintptr, args ...uintptr) int {
	switch len(a) {
	case 0:
		return syscall.Syscall(fn, uintptr(len(a)), 0, 0, 0)
	case 1:
		return syscall.Syscall(fn, uintptr(len(a)), a[0], 0, 0)
	case 2:
		return syscall.Syscall(fn, uintptr(len(a)), a[0], a[1], 0)
	case 3:
		return syscall.Syscall(fn, uintptr(len(a)), a[0], a[1], a[2])
	case 4:
		return syscall.Syscall6(fn, uintptr(len(a)), a[0], a[1], a[2], a[3], 0, 0)
	case 5:
		return syscall.Syscall6(fn, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], 0)
	case 6:
		return syscall.Syscall6(fn, uintptr(len(a)), a[0], a[1], a[2], a[3], a[4], a[5])
	default:
		panic("call() with too many arguments")
	}
	return -1
}
