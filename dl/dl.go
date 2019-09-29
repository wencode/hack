package dl

type Lib uintptr

type DLError struct {
	filename string
	errno uintptr
	errstr string
}

func (e DLError) Error() string {
	return e.filename + ": " + e.errstr
}
