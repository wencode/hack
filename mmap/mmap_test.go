package mmap

import (
	"fmt"
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

func BenchmarkCopyBuffer(b *testing.B) {
	for i := 0; i < 1000; i++ {
		fmt.Fprintf(os.Stderr, "benchmark %d-%d\n", b.N, i)
		doCopyBuffer(b)
	}
}

func doCopyBuffer(b *testing.B) {
	filename := fmt.Sprintf("mem_%d.data", b.N)
	siz := roundTo4096(b.N)
	if siz > 1024*1024*512 {
		siz = 1024 * 1024 * 512
	}
	mf, err := Open(
		filename,
		WithWrite(),
		WithLength(siz),
	)
	if err != nil {
		b.Fatalf("open %s error:%v", filename, err)
	}
	defer mf.Close()

	b.StartTimer()

	buf := mf.Buffer()
	if buf == nil {
		b.Fatalf("get %s buffer failed", filename)
	}

	pos := len(buf) / 2
	copy(buf[:pos], buf[pos:])

	b.StopTimer()
}

func roundTo4096(n int) int {
	return (n + 4095) & ^4095
}
