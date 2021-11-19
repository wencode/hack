package mmap

import (
	"os"
	"testing"
)

func TestMap(t *testing.T) {
	mf, err := Open("foo.data", WithWrite())
	if err != nil {
		t.Fatalf("open error %v", err)
	}
	defer mf.Close()

	buf := mf.Buffer()
	if len(buf) != os.Getpagesize() {
		t.Fatal("file default size error")
	}

	for i := 0; i < len(buf); i++ {
		buf[i] = '0' + byte(i%10)
	}
	buf.Sync()
}

func TestWrite(t *testing.T) {
	mf, err := Open("foo.data", WithWrite(), WithTruncate(), WithLength(4096))
	if err != nil {
		t.Fatalf("open error %v", err)
	}
	defer mf.Close()

	val := []byte{'a', 'b', 'c', 'd'}
	n, err := mf.Write(val)
	if err != nil || n != len(val) {
		t.Fatalf("write %d, err %v", n, err)
	}
	mf.Flush(mf.Buffer()[:4])
}

func TestRead(t *testing.T) {
	mf, err := Open("foo.data")
	if err != nil {
		t.Fatalf("open error %v", err)
	}
	defer mf.Close()

	val := make([]byte, 4)
	n, err := mf.Read(val)
	if err != nil || n != len(val) {
		t.Fatalf("write %d, err %v", n, err)
	}
}

func TestResize(t *testing.T) {
	mf, err := Open("foo.data", WithWrite(), WithTruncate(), WithLength(4096))
	if err != nil {
		t.Fatalf("open error %v", err)
	}
	defer mf.Close()

	if err := mf.Resize(8192, WithLength(8192)); err != nil {
		t.Fatalf("resize error %v", err)
	}
	data := []byte("hello")
	n := copy(mf.master[8185:], data)
	if n != len(data) {
		t.Errorf("write data error after resize")
	}
}

func TestExtend(t *testing.T) {
	mf, err := Open("foo.data", WithWrite(), WithTruncate(), WithLength(4096))
	if err != nil {
		t.Fatalf("open error %v", err)
	}
	defer mf.Close()

	buf, err := mf.ExtendMap(4096, 4096)
	if err != nil {
		t.Fatalf("extend error %v", err)
	}
	data := []byte("hello")
	n := copy(buf, data)
	if n != len(data) {
		t.Errorf("write data error after extend")
	}
}
