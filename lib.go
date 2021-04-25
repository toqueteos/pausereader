// Package pausereader provides pausable features into any io.Reader.
package pausereader

import (
	"io"
	"sync"
	"sync/atomic"
)

const (
	statePaused = uint32(0)
	stateActive = uint32(1)
)

// PausableReader defines a representation of pausable `io.Reader`s.
// It is implemented by PauseReader.
type PausableReader interface {
	io.Reader
	Pause()
	Resume()
	IsPaused() bool
}

// PauseReader wraps an `io.Reader` and allows to effectively pause all `Read`
// calls to it, turning it into a blocking operation.
type PauseReader struct {
	wg    waitGroup
	r     io.Reader
	state uint32
}

// New returns a PauseReader wrapping `r`.
func New(r io.Reader) *PauseReader {
	return &PauseReader{
		wg:    &stdWaitGroup{},
		r:     r,
		state: 1,
	}
}

// Pause pauses all future calls to `Read` making them blocking.
func (r *PauseReader) Pause() {
	oldValue := stateActive
	newValue := statePaused
	if atomic.CompareAndSwapUint32(&r.state, oldValue, newValue) {
		r.wg.Add(1)
	}
}

// Resume resumes all pending calls to `Read`, if the underlying `io.Reader` is
// not goroutine-safe then expect bad things to happen.
func (r *PauseReader) Resume() {
	oldValue := statePaused
	newValue := stateActive
	if atomic.CompareAndSwapUint32(&r.state, oldValue, newValue) {
		r.wg.Done()
	}
}

// IsPaused returns true if this PauseReader is paused or false otherwise.
func (r *PauseReader) IsPaused() bool {
	return atomic.LoadUint32(&r.state) == statePaused
}

// Read implements io.Reader it'll block if Pause is called until Resume is called.
func (r *PauseReader) Read(buf []byte) (n int, err error) {
	if r.IsPaused() {
		r.wg.Wait()
	}

	return r.r.Read(buf)
}

type waitGroup interface {
	Add(delta int)
	Done()
	Wait()
}

type stdWaitGroup struct {
	wg sync.WaitGroup
}

func (wg *stdWaitGroup) Add(delta int) { wg.wg.Add(delta) }
func (wg *stdWaitGroup) Done()         { wg.wg.Done() }
func (wg *stdWaitGroup) Wait()         { wg.wg.Wait() }
