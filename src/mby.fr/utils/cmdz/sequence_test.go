package cmdz

import (
	//"fmt"
	//"io"
	//"context"
	//"log"
	//"os/exec"
	"strings"
	"testing"
	"time"

	//"mby.fr/utils/promise"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	e1 = Cmd("echo", "foo")
	e2 = Cmd("echo", "bar")
	e3 = Cmd("echo", "baz")

	sleep10ms  = Cmd("sleep", "0.01")
	sleep11ms  = Cmd("sleep", "0.011")
	sleep100ms = Cmd("sleep", "0.1")
	sleep200ms = Cmd("sleep", "0.2")

	f1 = Cmd("false")
	f2 = Cmd("false")
)

func TestSerial(t *testing.T) {
	e1.reset()
	e2.reset()
	e3.reset()

	s := Serial(e1)
	assert.Equal(t, "echo foo", s.String())

	rc1, err1 := s.BlockRun()
	require.NoError(t, err1)
	assert.Equal(t, 0, rc1)
	assert.Equal(t, "foo\n", e1.StdoutRecord())
	assert.Equal(t, "", e2.StdoutRecord())
	assert.Equal(t, "", e3.StdoutRecord())
	assert.Equal(t, "foo\n", s.StdoutRecord())
	assert.Equal(t, "", s.StderrRecord())

	s2 := Serial(e1, e2)
	assert.Equal(t, "echo foo\necho bar", s2.String())

	rc2, err2 := s2.BlockRun()
	require.NoError(t, err2)
	assert.Equal(t, 0, rc2)
	assert.Equal(t, "foo\n", e1.StdoutRecord())
	assert.Equal(t, "bar\n", e2.StdoutRecord())
	assert.Equal(t, "", e3.StdoutRecord())
	assert.Equal(t, "foo\nbar\n", s2.StdoutRecord())
	assert.Equal(t, "", s2.StderrRecord())

	s2.Add(e3)
	assert.Equal(t, "echo foo\necho bar\necho baz", s2.String())

	rc3, err3 := s2.BlockRun()
	require.NoError(t, err3)
	assert.Equal(t, 0, rc3)
	assert.Equal(t, "foo\n", e1.StdoutRecord())
	assert.Equal(t, "bar\n", e2.StdoutRecord())
	assert.Equal(t, "baz\n", e3.StdoutRecord())
	assert.Equal(t, "foo\nbar\nbaz\n", s2.StdoutRecord())
	assert.Equal(t, "", s2.StderrRecord())

	// Test serial timings
	s2.Add(sleep10ms)
	start := time.Now()
	_, err := s2.BlockRun()
	duration := time.Since(start).Microseconds()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, duration, int64(10000), "Serial too quick !")
	assert.Less(t, duration, int64(20000), "Serial too slow !")
	assert.Equal(t, "foo\nbar\nbaz\n", s2.StdoutRecord())
	assert.Equal(t, "", s2.StderrRecord())

	s2.Add(sleep10ms)
	start = time.Now()
	_, err = s2.BlockRun()
	duration = time.Since(start).Microseconds()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, duration, int64(20000), "Serial too quick !")
	assert.Less(t, duration, int64(40000), "Serial too slow !")
	assert.Equal(t, "foo\nbar\nbaz\n", s2.StdoutRecord())
	assert.Equal(t, "", s2.StderrRecord())

	s2.Add(sleep100ms)
	start = time.Now()
	_, err = s2.BlockRun()
	duration = time.Since(start).Microseconds()
	require.NoError(t, err)
	assert.GreaterOrEqual(t, duration, int64(120000), "Serial too quick !")
	assert.Less(t, duration, int64(150000), "Serial too slow !")
	assert.Equal(t, "foo\nbar\nbaz\n", s2.StdoutRecord())
	assert.Equal(t, "", s2.StderrRecord())
}

func TestSerial_Retries(t *testing.T) {
	e1.reset()
	f1.reset()
	e2.reset()
	s := Serial(e1, f1, e2)
	rc, err := s.BlockRun()

	require.NoError(t, err)
	assert.Equal(t, 0, rc)
	assert.Equal(t, []int{0}, e1.ResultsCodes)
	assert.Equal(t, []int{1}, f1.ResultsCodes)
	assert.Equal(t, []int{0}, e2.ResultsCodes)

	s.Retries(2, 10)
	rc, err = s.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 0, rc)
	assert.Equal(t, []int{0}, e1.ResultsCodes)
	assert.Equal(t, []int{1, 1, 1}, f1.ResultsCodes)
	assert.Equal(t, []int{0}, e2.ResultsCodes)

	s.Add(f1)
	rc, err = s.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
}

func TestSerial_Outputs(t *testing.T) {
	e1.reset()
	f1.reset()
	e2.reset()

	sb := strings.Builder{}
	s := Serial(e1, f1, e2).Outputs(&sb, nil)
	rc, err := s.BlockRun()

	require.NoError(t, err)
	assert.Equal(t, 0, rc)
	assert.Equal(t, "foo\nbar\n", sb.String())
}

func TestSerial_FailFast(t *testing.T) {
	e1.reset()
	f1.reset()
	e2.reset()

	s := Serial(e1, f1, e2).FailFast(true)
	rc, err := s.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, []int{0}, e1.ResultsCodes)
	assert.Equal(t, []int{1}, f1.ResultsCodes)
	assert.Equal(t, []int(nil), e2.ResultsCodes)
	assert.Equal(t, "foo\n", s.StdoutRecord())

	s.Retries(2, 10)
	rc, err = s.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, []int{0}, e1.ResultsCodes)
	assert.Equal(t, []int{1, 1, 1}, f1.ResultsCodes)
	assert.Equal(t, []int(nil), e2.ResultsCodes)
	assert.Equal(t, "foo\n", s.StdoutRecord())
}


func TestParallel(t *testing.T) {
	e1.reset()
	e2.reset()
	e3.reset()

	p := Parallel(e1)
	assert.Equal(t, "echo foo", p.String())

	rc, err := p.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 0, rc)
	assert.Equal(t, "foo\n", e1.StdoutRecord())
	assert.Equal(t, "", e2.StdoutRecord())
	assert.Equal(t, "", e3.StdoutRecord())
	assert.Equal(t, "foo\n", p.StdoutRecord())
	assert.Equal(t, "", p.StderrRecord())

	p2 := Parallel(e1, e2)
	assert.Equal(t, "echo foo\necho bar", p2.String())

	rc2, err2 := p2.BlockRun()
	require.NoError(t, err2)
	assert.Equal(t, 0, rc2)
	assert.Equal(t, "foo\n", e1.StdoutRecord())
	assert.Equal(t, "bar\n", e2.StdoutRecord())
	assert.Equal(t, "", e3.StdoutRecord())
	assert.Equal(t, "foo\nbar\n", p2.StdoutRecord())
	assert.Equal(t, "", p2.StderrRecord())

	p2.Add(e3)
	assert.Equal(t, "echo foo\necho bar\necho baz", p2.String())

	rc3, err3 := p2.BlockRun()
	require.NoError(t, err3)
	assert.Equal(t, 0, rc3)
	assert.Equal(t, "foo\n", e1.StdoutRecord())
	assert.Equal(t, "bar\n", e2.StdoutRecord())
	assert.Equal(t, "baz\n", e3.StdoutRecord())
	assert.Equal(t, "foo\nbar\nbaz\n", p2.StdoutRecord())
	assert.Equal(t, "", p2.StderrRecord())

	// Test serial timings
	p2.Add(sleep10ms)
	start := time.Now()
	_, err = p2.BlockRun()
	duration := time.Since(start).Microseconds()
	require.NoError(t, err)
	//assert.Len(t, 4, rc)
	assert.GreaterOrEqual(t, duration, int64(10000), "Parallel too quick !")
	assert.Less(t, duration, int64(20000), "Parallel too slow !")
	assert.Equal(t, "foo\nbar\nbaz\n", p2.StdoutRecord())
	assert.Equal(t, "", p2.StderrRecord())

	p2.Add(sleep11ms)
	start = time.Now()
	_, err = p2.BlockRun()
	duration = time.Since(start).Microseconds()
	require.NoError(t, err)
	//assert.Len(t, 5, rc)
	assert.GreaterOrEqual(t, duration, int64(11000), "Parallel too quick !")
	assert.Less(t, duration, int64(21000), "Parallel too slow !")
	assert.Equal(t, "foo\nbar\nbaz\n", p2.StdoutRecord())
	assert.Equal(t, "", p2.StderrRecord())

	p2.Add(sleep100ms)
	start = time.Now()
	_, err = p2.BlockRun()
	duration = time.Since(start).Microseconds()
	require.NoError(t, err)
	//assert.Len(t, 5, rc)
	assert.GreaterOrEqual(t, duration, int64(100000), "Parallel too quick !")
	assert.Less(t, duration, int64(120000), "Parallel too slow !")
	assert.Equal(t, "foo\nbar\nbaz\n", p2.StdoutRecord())
	assert.Equal(t, "", p2.StderrRecord())
}


func TestParallel_Retries(t *testing.T) {
	e1.reset()
	f1.reset()
	e2.reset()
	p := Parallel(e1, f1, e2)
	rc, err := p.BlockRun()

	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, []int{0}, e1.ResultsCodes)
	assert.Equal(t, []int{1}, f1.ResultsCodes)
	assert.Equal(t, []int{0}, e2.ResultsCodes)

	p.Retries(2, 10)
	rc, err = p.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, []int{0}, e1.ResultsCodes)
	assert.Equal(t, []int{1, 1, 1}, f1.ResultsCodes)
	assert.Equal(t, []int{0}, e2.ResultsCodes)

	p.Add(f2)
	rc, err = p.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
}

func TestParallel_Outputs(t *testing.T) {
	c1 := Cmd("/bin/sh", "-c", "sleep 0.1 ; echo foo")
	c2 := Cmd("/bin/sh", "-c", "sleep 0.2 ; echo bar")
	f1.reset()

	sb := strings.Builder{}
	p := Parallel(c1, f1, c2).Outputs(&sb, nil)
	rc, err := p.BlockRun()

	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, "foo\nbar\n", sb.String())
}

func TestParallel_FailFast(t *testing.T) {
	c1 := Cmd("/bin/sh", "-c", "sleep 0.1 && echo foo")
	c2 := Cmd("/bin/sh", "-c", "sleep 0.2 && echo bar")
	f1.reset()

	p := Parallel(c1, f1, c2).FailFast(true)
	rc, err := p.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, "", p.StdoutRecord())
	assert.Equal(t, []int(nil), c1.ResultsCodes)
	assert.Equal(t, []int{1}, f1.ResultsCodes)
	assert.Equal(t, []int(nil), c2.ResultsCodes)
	
	c1.reset()
	c2.reset()
	f1.reset()
	p.Retries(2, 10)
	rc, err = p.BlockRun()
	require.NoError(t, err)
	assert.Equal(t, 1, rc)
	assert.Equal(t, []int(nil), c1.ResultsCodes)
	assert.Equal(t, []int{1, 1, 1}, f1.ResultsCodes)
	assert.Equal(t, []int(nil), c2.ResultsCodes)
	assert.Equal(t, "", p.StdoutRecord())
}
