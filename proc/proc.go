package proc

type Proc uintptr

func (p Proc) Call(args ...uintptr) uintptr {
	return Call(uintptr(p), args...)
}

