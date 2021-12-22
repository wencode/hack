// +build !amd64

package mmap

import (
	"os"
	"syscall"
)

var (
	pageSize     uintptr = uintptr(os.Getpagesize())
	pageSizeMask uintptr = ^(pageSize - 1)
)

func Flush(addr, length uintptr) uintptr {
	for pageAddr := uintptr(addr & pageSizeMask); pageAddr < addr+length; pageAddr += pageSize {
		_, _, err := syscall.Syscall(syscall.SYS_MSYNC, pageAddr, pageSize, syscall.MS_SYNC)
		if err != 0 {
			return uintptr(err)
		}
	}
	return 0
}
