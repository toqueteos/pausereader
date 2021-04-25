package pausereader

import (
	"io"
	"strings"
	"testing"

	"github.com/matryer/is"
)

func reader() (*PauseReader, *testWaitGroup) {
	r := New(strings.NewReader("pausereader"))
	wg := wg()
	r.wg = wg
	return r, wg
}

func read(r io.Reader) (n int, err error) {
	buf := make([]byte, 4)
	n, err = r.Read(buf)
	return
}

func TestPauseReader(t *testing.T) {
	is := is.New(t)

	r, wg := reader()

	r.Pause()
	is.Equal(wg.add, 1)          // wg.Add is called once
	is.Equal(r.IsPaused(), true) // state is paused

	read(r)
	is.Equal(wg.wait, 1) // wg.Wait called

	r.Resume()
	is.Equal(wg.done, 1)          // wg.Done is called once
	is.Equal(r.IsPaused(), false) // state is active

	read(r)
	is.Equal(wg.wait, 1) // wg.Wait not called
}

func TestResume(t *testing.T) {
	is := is.New(t)

	r, wg := reader()

	is.Equal(r.IsPaused(), false) // default state is active

	r.Resume()
	r.Resume()
	r.Resume()

	is.Equal(wg.add, 0)           // wg.Add not called
	is.Equal(r.IsPaused(), false) // state is active

	read(r)
	is.Equal(wg.wait, 0) // wg.Wait not called
}

func TestPause(t *testing.T) {
	is := is.New(t)

	r, wg := reader()

	r.Pause()
	r.Pause()
	r.Pause()

	is.Equal(wg.add, 1)          // wg.Add called once
	is.Equal(r.IsPaused(), true) // state is paused

	read(r)
	is.Equal(wg.wait, 1) // wg.Wait called once
}

type testWaitGroup struct {
	add  int
	done int
	wait int
}

func wg() *testWaitGroup {
	return &testWaitGroup{}
}

func (twg *testWaitGroup) Add(delta int) {
	twg.add++
}
func (twg *testWaitGroup) Done() {
	twg.done++
}
func (twg *testWaitGroup) Wait() {
	twg.wait++
}
