package dl

import (
	"syscall"
)

// Implemented in runtime/syscall_windows.goc; we provide jumps to them in our assembly file.
func loadlibrary(filename *uint16) (handle uintptr, err syscall.Errno)
func getprocaddress(handle uintptr, procname *uint8) (proc uintptr, err syscall.Errno)

func Open(filename string) (lib Lib, err error) {
	utf16filename, e := syscall.UTF16PtrFromString(filename)
	if e != nil {
		err = &DLError{
			filename: filename,
			errstr:   "convert to utf16 failed: " + e.Error(),
		}
		return
	}
	handle, e := loadlibrary(utf16filename)
	if handle == 0 {
		err = &DLError{
			filename: filename,
			errstr:   "load " + filename + "filed: " + e.Error(),
		}
		return
	}
	return Lib(handle), nil

}

func (lib Lib) Close() {
	if closeHandle == 0 {
		return
	}
	syscall.Syscall(closeHandle, 1, uintptr(lib), 0, 0)
}

func (lib Lib) Sym(symbol string) uintptr {
	utf8symbol := uint8ptr(symbol)
	p, e := getprocaddress(uintptr(lib), utf8symbol)
	if e != 0 {
		return 0
	}
	return p
}

func uint8ptr(s string) *uint8 {
	b := make([]byte, len(s)+1)
	copy(b, s)
	return &b[0]
}

var (
	knlLib      Lib
	closeHandle uintptr
)

func init() {
	knlLib, e := Open("kernel32.dll")
	if e != nil {
		return
	}
	closeHandle = knlLib.Sym("CloseHandle")
}
