package mmap

import (
	"errors"
	"io"
	"os"
	"reflect"
	"syscall"
	"unsafe"
)

const (
	RDONLY = 0
	RDWR   = 1 << iota
	COW
	EXEC
)

var (
	ErrArgument = errors.New("argument error")
)

type MapBuf []byte

type MapFile struct {
	file   *os.File
	prot   int
	offset int
	buf    MapBuf
	index  int
}

type Param struct {
	Offset   int
	Len      int
	Prot     int
	Private  bool
	Append   bool
	Truncate bool
}

type Option func(param *Param)

func WithOffset(offset int) Option {
	return func(param *Param) {
		param.Offset = offset
	}
}

func WithLength(len int) Option {
	return func(param *Param) {
		param.Len = len
	}
}

func WithWrite() Option {
	return func(param *Param) {
		param.Prot |= RDWR
	}
}

func WithCopyOnWrite() Option {
	return func(param *Param) {
		param.Prot |= COW
	}
}

func WithPrivate() Option {
	return func(param *Param) {
		param.Private = true
	}
}

func WithTruncate() Option {
	return func(param *Param) {
		param.Truncate = true
	}
}

func parseMapParam(opts ...Option) *Param {
	param := &Param{}
	for _, opt := range opts {
		opt(param)
	}
	return param
}

func Open(filename string, opts ...Option) (*MapFile, error) {
	param := parseMapParam(opts...)

	var (
		file *os.File
		err  error
	)
	if filename != "" {
		file, err = OpenFile(filename, param.Prot, param.Truncate)
		if err != nil {
			return nil, err
		}
	}

	if file == nil && param.Len == 0 {
		return nil, ErrArgument
	}

	return _open(file, param)
}

func OpenWithFile(file *os.File, opts ...Option) (*MapFile, error) {
	param := parseMapParam(opts...)
	return _open(file, param)
}

func _open(file *os.File, param *Param) (*MapFile, error) {
	if err := checkFile(file, param); err != nil {
		return nil, err
	}

	var fd = -1
	if file != nil {
		fd = int(file.Fd())
	}
	buf, err := Mmap(fd, param.Prot, param.Offset, param.Len)
	if err != nil {
		return nil, err
	}

	return &MapFile{
		file:   file,
		prot:   param.Prot,
		offset: param.Offset,
		buf:    buf,
	}, nil
}

func (m *MapFile) Close() {
	m.Unmap()
	if m.file != nil {
		m.file.Close()
		m.file = nil
	}
}

func (m *MapFile) Unmap() {
	if m.buf != nil {
		m.buf.Unmap()
		m.buf = nil
	}
}

func (m *MapFile) Remap(offset, len int) error {
	m.Unmap()
	buf, err := Mmap(int(m.file.Fd()), m.prot, offset, len)
	if err != nil {
		return err
	}
	m.buf = buf
	m.index = 0
	return nil
}

func (m *MapFile) Buffer() MapBuf {
	return m.buf
}

func (m *MapFile) Read(b []byte) (n int, err error) {
	if m.buf == nil || m.index >= len(m.buf) {
		err = io.EOF
		return
	}
	src := m.buf[m.index:]
	n = copy(b, src)
	m.index += n
	return
}

func (m *MapFile) Write(b []byte) (n int, err error) {
	if m.buf == nil || m.index >= len(m.buf) {
		err = io.EOF
		return
	}
	dst := m.buf[m.index:]
	n = copy(dst, b)
	m.index += n
	return
}

func (m *MapFile) Seek(offset int64, whence int) (newoff int64, err error) {
	switch whence {
	case io.SeekStart:
		if offset < 0 || offset > int64(len(m.buf)) {
			err = ErrArgument
			return
		}
		m.index = int(offset)
	case io.SeekCurrent:
		cur := m.index + int(offset)
		if cur < 0 || cur > len(m.buf) {
			err = ErrArgument
			return
		}
		m.index = cur
	case io.SeekEnd:
		if offset < 0 || offset > int64(len(m.buf)) {
			err = ErrArgument
			return
		}
		m.index = len(m.buf) - 1 - int(offset)
	default:
		err = ErrArgument
		return
	}
	newoff = int64(m.index)
	return
}

func (m *MapFile) Sync() {
	if m.buf == nil {
		return
	}
	m.buf.Sync()
}

func (m *MapFile) Flush(b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	baseh := m.buf.header()
	if bh.Data < baseh.Data || bh.Data >= baseh.Data+uintptr(baseh.Len) {
		return
	}
	Flush(bh.Data, uintptr(bh.Len))
}

func (m *MapFile) Resize(newSize int, opts ...Option) error {
	param := parseMapParam(opts...)
	if param.Len > 0 {
		if param.Offset < 0 {
			return ErrArgument
		}
		if param.Offset+param.Len > newSize || param.Offset >= newSize {
			return ErrArgument
		}
	}
	if m.file == nil {
		return ErrArgument
	} else {
		st, err := m.file.Stat()
		if err != nil {
			return err
		}
		if int64(newSize) <= st.Size() {
			return ErrArgument
		}
	}

	fillFile(m.file, newSize)
	if param.Len > 0 {
		if err := m.Remap(param.Offset, param.Len); err != nil {
			return err
		}
	}
	return nil
}

func OpenFile(filename string, prot int, truncate bool) (*os.File, error) {
	flags := os.O_RDONLY
	if prot != RDONLY {
		flags = os.O_RDWR
	}
	flags |= os.O_CREATE
	if truncate {
		flags |= os.O_TRUNC
	}
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	file, err := os.OpenFile(filename, flags, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func checkFile(file *os.File, param *Param) error {
	st, err := file.Stat()
	if err != nil {
		return err
	}

	filesize := int(st.Size())
	fill := false
	if param.Len == 0 {
		if ps := os.Getpagesize(); filesize < ps {
			filesize = ps
			fill = true
		}
		param.Len = filesize
	} else {
		if filesize < param.Len {
			filesize = param.Len
			fill = true
		}
	}

	if fill {
		fillFile(file, param.Len)
	}

	return nil
}

func fillFile(file *os.File, length int) {
	var (
		tmp [1]byte
	)
	file.Seek(int64(length-1), 0)
	file.Write(tmp[:])
	file.Sync()
	file.Seek(0, 0)
}

func (mb MapBuf) header() *reflect.SliceHeader {
	return (*reflect.SliceHeader)(unsafe.Pointer(&mb))
}
