package pausereader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type PauseSuite struct{}

var _ = Suite(&PauseSuite{})

func (s *PauseSuite) SetUpSuite(c *C) {}

func (s *PauseSuite) TestInMemoryReader(c *C) {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, "Hello, playground")

	var pr = NewPauseReader(&buf, 100, 2)

	pr.Pause()
	c.Log(".Pause")
	go func() {
		t := time.Now()
		time.Sleep(5 * time.Second)
		c.Log(".Resume")
		c.Logf("backoff=%d", pr.backoff)
		pr.Resume()
		c.Logf("elapsed=%s", time.Since(t))
	}()

	var out bytes.Buffer
	_, err := io.CopyN(&out, pr, 20)
	c.Assert(err, Equals, io.EOF)

	c.Logf("buf=%q", out.String())
	c.Logf("spins=%d", pr.spins)
	c.Assert(pr.spins <= 200, Equals, true)
}

func (s *PauseSuite) TestHttpReader(c *C) {
	res, err := http.Get("https://github.com/toqueteos/webbrowser/archive/master.zip")
	c.Assert(err, IsNil)
	defer res.Body.Close()

	var pr = NewPauseReader(res.Body, 10, 2)

	pr.Pause()
	c.Log(".Pause")
	go func() {
		t := time.Now()
		time.Sleep(2 * time.Second)
		c.Log(".Resume")
		c.Logf("backoff=%d", pr.backoff)
		pr.Resume()
		c.Logf("elapsed=%s", time.Since(t))

		t = time.Now()
		pr.Pause()
		time.Sleep(2 * time.Second)
		c.Log(".Resume")
		c.Logf("backoff=%d", pr.backoff)
		pr.Resume()
		c.Logf("elapsed=%s", time.Since(t))
	}()

	var out bytes.Buffer
	_, err = io.Copy(&out, pr)
	c.Assert(err, IsNil)

	c.Logf("buf=%q", out.Len())
	c.Logf("spins=%d", pr.spins)
	c.Assert(pr.spins <= 200, Equals, true)
}
