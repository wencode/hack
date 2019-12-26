// +build darwin linux

package mmap

import "golang.org/x/sys/unix"

func Mmap(fd, prot, offset, len int) (MapBuf, error) {
	var (
		flags  = unix.MAP_SHARED
		s_prot = unix.PROT_READ
	)
	if prot&RDWR != 0 {
		s_prot |= unix.PROT_WRITE
	}
	if prot&COW != 0 {
		s_prot |= unix.PROT_WRITE
		flags = unix.MAP_PRIVATE
	}
	if prot&EXEC != 0 {
		s_prot |= unix.PROT_EXEC
	}
	if fd < 0 {
		// The mapping is not backed by any file;
		// it contents are initialize to zero.
		// the fd and offset arguments are ignored.
		flags |= unix.MAP_ANON
	}

	buf, err := unix.Mmap(fd, int64(offset), len, s_prot, flags)
	if err != nil {
		return nil, err
	}
	return MapBuf(buf), nil
}

func (mb MapBuf) Unmap() error {
	return unix.Munmap([]byte(mb))
}

func (mb MapBuf) Mlock() error {
	return unix.Mlock([]byte(mb))
}

func (mb MapBuf) Munlock() error {
	return unix.Munlock([]byte(mb))
}

func (mb MapBuf) Sync() error {
	return unix.Msync([]byte(mb), unix.MS_ASYNC)
}
