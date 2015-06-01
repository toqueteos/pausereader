package pausereader

import (
	"io"
	"sync"
	"time"
)

type Resumer interface {
	Pause()
	Resume()
	Running() bool
}

// PauseReader wraps an `io.Reader` and allows to effectively pause all .Read
// calls to it, basically returning `(0, nil)`.
type PauseReader struct {
	sync.RWMutex
	r       io.Reader
	state   bool
	base    int // Base wait time in milliseconds
	checks  int // Number of sleeps to increase backoff
	c       int
	backoff int // Exponential backoff
	spins   int
}

// NewPauseReader returns a PauseReader wrapping `r` with `wait` milliseconds
// of base wait time, `wait` will exponentially increase every `c` checks
// PauseReader is still paused.
func NewPauseReader(r io.Reader, wait, c int) *PauseReader {
	return &PauseReader{
		r:       r,
		state:   false,
		base:    wait,
		checks:  c,
		backoff: 1,
	}
}
func (p *PauseReader) Pause()        { p.Lock(); defer p.Unlock(); p.state = false }
func (p *PauseReader) Resume()       { p.Lock(); defer p.Unlock(); p.c = 0; p.backoff = 1; p.state = true }
func (p *PauseReader) Running() bool { return p.state }
func (p *PauseReader) Read(buf []byte) (n int, err error) {
	p.RLock()
	var state = p.state
	var t = p.base * p.backoff
	p.RUnlock()
	if state {
		return p.r.Read(buf)
	} else {
		p.Lock()
		p.spins++
		if p.c < p.checks {
			p.c++
		} else {
			p.c = 0
			p.backoff *= 2
		}
		p.Unlock()
		time.Sleep(time.Duration(t) * time.Millisecond)
		buf = nil
		return 0, nil
	}
}
