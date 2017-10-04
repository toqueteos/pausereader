package pausereader

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type oneByteReader struct {
	buf []byte
	pos int
	mtx sync.Mutex
}

func (r *oneByteReader) Read(buf []byte) (n int, err error) {
	r.mtx.Lock()
	buf[0] = r.buf[r.pos]
	r.pos++
	r.mtx.Unlock()
	return 1, nil
}

func TestInMemoryReader(t *testing.T) {
	buf := &oneByteReader{buf: []byte{1, 2, 3, 4, 5, 6, 7, 8}}

	r := New(buf)
	r.Pause()

	counter := uint32(0)

	for i := 0; i < len(buf.buf); i++ {
		go read(t, r, &counter)
	}

	time.Sleep(200 * time.Millisecond)
	if counter != 0 {
		t.Fatalf("wrong state: expected counter=0, got counter=%d", counter)
	}

	r.Resume()
	time.Sleep(200 * time.Millisecond)
	if r.IsPaused() {
		t.Fatal("wrong PausableReader state: expected paused=false, got paused=true")
	}

	expectedCounter := uint32(len(buf.buf))
	gotCounter := atomic.LoadUint32(&counter)
	if gotCounter != expectedCounter {
		t.Fatalf("wrong state: expected counter=%d, got counter=%d", expectedCounter, gotCounter)
	}
}

func read(t *testing.T, r PausableReader, counter *uint32) {
	buf := make([]byte, 1)
	n, err := r.Read(buf)
	if buf[0] == 0 {
		t.Fatal("wrong state: expected value != 0")
	}
	if n != 1 {
		t.Fatalf("wrong state: expected n=1, got n=%d", n)
	}
	if err != nil {
		t.Fatalf("wrong state: expected error=nil, got error=%v", err)
	}
	atomic.AddUint32(counter, 1)
}
