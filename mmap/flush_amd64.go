// +build amd64

package mmap

func flush_cache(addr, size uintptr)

const (
	cacheLineLength uintptr = 64
	cacheLineMask   uintptr = ^(cacheLineLength - 1)
)

func Flush(addr, length uintptr) uintptr {
	alignAddr := addr & cacheLineMask
	length += (addr - alignAddr)
	flush_cache(alignAddr, length)
	return 0
}
