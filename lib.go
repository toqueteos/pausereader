package pausereader

import (
	"io"
	"sync"
	"sync/atomic"
)

type PausableReader interface {
	io.Reader
	Pause()
	Resume()
	IsPaused() bool
}

// PauseReader wraps an `io.Reader` and allows to effectively pause all .Read
// calls to it, turn it a blocking operation.
type PauseReader struct {
	wg    sync.WaitGroup
	r     io.Reader
	state uint32
}

// New returns a PauseReader wrapping `r`.
func New(r io.Reader) *PauseReader {
	return &PauseReader{
		r:     r,
		state: 1,
	}
}
func (r *PauseReader) Pause() {
	if atomic.CompareAndSwapUint32(&r.state, 1, 0) {
		r.wg.Add(1)
	}
}

func (r *PauseReader) Resume() {
	if atomic.CompareAndSwapUint32(&r.state, 0, 1) {
		r.wg.Done()
	}
}

func (r *PauseReader) IsPaused() bool {
	return atomic.LoadUint32(&r.state) == 0
}

func (r *PauseReader) Read(buf []byte) (n int, err error) {
	if r.IsPaused() {
		r.wg.Wait()
	}

	return r.r.Read(buf)
}
