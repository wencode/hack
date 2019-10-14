// +build linux darwin

package dl

// #include <dlfcn.h>
// #include <stdlib.h>
// #include <string.h>
// #cgo LDFLAGS: -ldl -pthread
import "C"
import (
	"unsafe"
)

func Open(filename string) (lib Lib, err error) {
	s := C.CString(filename)
	defer C.free(unsafe.Pointer(s))

	handle, e := C.dlopen(s, C.RTLD_NOW)
	if handle == nil {
		err = newDLError(filename, e)
		return
	}
	lib = Lib(uintptr(handle))
	return
}

func (lib Lib) Close() {
	C.dlclose(unsafe.Pointer(lib))
}

func (lib Lib) Sym(symbol string) uintptr {
	s := C.CString(symbol)
	defer C.free(unsafe.Pointer(s))

	addr := C.dlsym(unsafe.Pointer(lib), s)
	return uintptr(addr)
}

func newDLError(filename string, err error) *DLError {
	return &DLError{
		filename:filename,
		errstr: err.Error(),
	}
}