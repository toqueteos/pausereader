// Package pausereader provides pausable features into any io.Reader.
package pausereader

import (
	"io"
	"sync"
	"sync/atomic"
)

// PausableReader defines a representation of pausable `io.Reader`s.
// It is implemented by PauseReader.
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

// Pause pauses all future calls to `Read` making them blocking.
func (r *PauseReader) Pause() {
	if atomic.CompareAndSwapUint32(&r.state, 1, 0) {
		r.wg.Add(1)
	}
}

// Resume resumes all pending calls to `Read`, if the underlying `io.Reader` is
// not goroutine-safe then expect bad things to happen.
func (r *PauseReader) Resume() {
	if atomic.CompareAndSwapUint32(&r.state, 0, 1) {
		r.wg.Done()
	}
}

// IsPaused returns true if this PauseReader is paused or false otherwise.
func (r *PauseReader) IsPaused() bool {
	return atomic.LoadUint32(&r.state) == 0
}

// Read implements io.Reader it'll block if Pause is called until Resume is called.
func (r *PauseReader) Read(buf []byte) (n int, err error) {
	if r.IsPaused() {
		r.wg.Wait()
	}

	return r.r.Read(buf)
}
