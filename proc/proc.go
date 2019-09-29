package proc

type Proc uintptr

func (p Proc) Call(args ...uintptr) int {
	return Call(uintptr(p), args...)
}

