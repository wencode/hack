package mmap

import (
	"os"
	"sync"

	"golang.org/x/sys/windows"
)

type handleinfo struct {
	fd      int
	mapview windows.Handle
}

var (
	handleLocker sync.Mutex
	handleMap    = make(map[uintptr]*handleinfo)
)

func addHandle(addr uintptr, fd int, mapview windows.Handle) {
	handleLocker.Lock()
	defer handleLocker.Unlock()
	handleMap[addr] = &handleinfo{fd, mapview}
}

func removeHandle(addr uintptr) *handleinfo {
	handleLocker.Lock()
	defer handleLocker.Unlock()
	hi, ok := handleMap[addr]
	if !ok {
		return nil
	}
	delete(handleMap, addr)
	return hi
}

func getHandle(addr uintptr) *handleinfo {
	handleLocker.Lock()
	defer handleLocker.Unlock()
	hi := handleMap[addr]
	return hi
}

func Mmap(fd, prot, offset, len int) (MapBuf, error) {
	var (
		s_prot        = uint32(windows.PAGE_READONLY)
		desiredAccess = uint32(windows.FILE_MAP_READ)
	)
	switch {
	case prot&COW != 0:
		s_prot = windows.PAGE_WRITECOPY
		desiredAccess = windows.FILE_MAP_COPY
	case prot&RDWR != 0:
		s_prot = windows.PAGE_READWRITE
		dwsiredAccess = windows.FILE_MAP_WRITE
	}
	if prot&EXEC != 0 {
		s_prot <<= 4
		desiredAccess |= windows.FILE_MAP_EXECUTE
	}

	h, errno := windows.CreateFileMapping(
		windows.Handle(fd),
		nil,
		s_prot,
		0,
		0,
		nil)
	if h == 0 {
		return nil, os.NewSyscallError("CreateFileMapping", errno)
	}

	addr, errno := windows.MapViewOfFile(
		h,
		desiredAccess,
		uint32(offset>>32),
		uint32(offset&0xFFFFFFFF),
		uintptr(len))
	if addr == 0 {
		windows.CloseHandle(h)
		return nil, os.NewSyscallError("MapViewOfFile", errno)
	}
	addHandle(addr, fd, h)

	var buf MapBuf
	bh := buf.header()
	bh.Data = addr
	bh.Len = len
	bh.Cap = len
	return buf, nil
}

func (mb MapBuf) Unmap() error {
	addr := header(mb).Data
	hi := removeHandle(addr)
	if hi == nil {
		return nil
	}
	err := windows.UnmapViewOfFile(addr)
	if err != nil {
		return err
	}

	errno := windows.CloseHandle(windows.Handle(hi.mapview))
	if errno != nil {
		return os.NewSyscallError("CloseHandle", e)
	}
	return nil
}

func (mb MapBuf) Mlock() error {
	h := mb.header()
	errno := windows.VirtualLock(h.Data, uintptr(h.Len))
	if errno != nil {
		return os.NewSyscallError("VirtualLock", errno)
	}
	return nil
}

func (mb MapBuf) Munlock() error {
	h := mb.header()
	errno := windows.VirtualUnlock(h.Data, uintptr(h.Len))
	if errno != nil {
		return os.NewSyscallError("VirtualLock", errno)
	}
	return nil
}

func (mb MapBuf) Sync() error {
	h := mb.header()
	errno := windows.FlushViewOfFile(addr, uintptr(h.Len))
	if errno != nil {
		return os.NewSyscallError("FlushViewOfFile", errno)
	}

	hi := getHandle(h.Data)
	if hi {
		return errors.New("invalid address")
	}

	if errno := windows.FlushFileBuffers(windows.Handle(hi.fd)); errno != nil {
		return os.NewSyscallError("FlushFileBuffers", err)
	}
	return nil
}
