package output

import (
	"io"
	"os"
	"bufio"
	"time"
	//"fmt"
)

// Outputs responsible for keeping reference of outputs writers (example: stdout, file, ...)

type Flusher interface {
	Flush() error
}

type Outputs interface {
	Flusher
	Out() io.Writer
	Err() io.Writer
	LastWriteTime() time.Time
}

type BasicOutputs struct {
	out, err *ActivityWriter
}

func (o BasicOutputs) Flush() error {
	return nil
}

func (o BasicOutputs) LastWriteTime() time.Time {
	if o.out.activity.Before(o.err.activity) {
		return o.err.activity
	}
	return o.out.activity
}

func (o BasicOutputs) Out() io.Writer {
	return o.out
}

func (o BasicOutputs) Err() io.Writer {
	return o.err
}

type BufferedOutputs struct {
	BasicOutputs
	bufferedOut, bufferedErr *bufio.Writer
}
func (o BufferedOutputs) Flush() error {
	outs := []io.Writer{o.bufferedOut, o.bufferedErr}
	for _, out := range outs {
		f, ok := out.(Flusher)
		if ok {
			err := f.Flush()
			if err != nil {
				return err
			}
		}
	}
	//fmt.Printf("flushed %T %T\n", o.out, o.err)
	return nil
}

type ActivityWriter struct {
	nested io.Writer
	activity time.Time
}
func (w *ActivityWriter) Write(b []byte) (int, error) {
	t := time.Now()
	w.activity = t
	return w.nested.Write(b)
}

func New(out, err io.Writer) BasicOutputs {
	t := time.Time{}
	activityOut := ActivityWriter{out, t}
	activityErr := ActivityWriter{err, t}
	return BasicOutputs{&activityOut, &activityErr}
}

func NewStandardOutputs() Outputs {
	return New(os.Stdout, os.Stderr)
}

func NewBufferedOutputs(outputs Outputs) Outputs {
	buffOut := bufio.NewWriter(outputs.Out())
	buffErr := bufio.NewWriter(outputs.Err())
	basic := New(buffOut, buffErr)
	buffered := BufferedOutputs{basic, buffOut, buffErr}
	return buffered
}

// TODO ?
//func NewFileOutputs(outFilePath, errFilePath string) Outputs {
//	if outFilePath == errFilePath || errFilePath == "" {
//		// Same file for both outputs
//
//	}
//}
