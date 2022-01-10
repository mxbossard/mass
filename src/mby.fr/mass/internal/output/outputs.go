package output

import (
	"io"
	"os"
	"bufio"
	"time"
	//"fmt"

	"mby.fr/utils/datetime"
)

// Outputs responsible for keeping reference of outputs writers (example: stdout, file, ...)

type Flusher interface {
	Flush() error
}

type Outputs interface {
	Flusher
	Log() io.Writer
	Out() io.Writer
	Err() io.Writer
	LastWriteTime() time.Time
}

type BasicOutputs struct {
	log, out, err *ActivityWriter
}

func (o BasicOutputs) Flush() error {
	return nil
}

func (o BasicOutputs) LastWriteTime() time.Time {
	lastTime := datetime.Max(o.log.activity, o.out.activity, o.err.activity)
	return lastTime
}

func (o BasicOutputs) Log() io.Writer {
	return o.log
}

func (o BasicOutputs) Out() io.Writer {
	return o.out
}

func (o BasicOutputs) Err() io.Writer {
	return o.err
}

type BufferedOutputs struct {
	BasicOutputs
	bufferedLog, bufferedOut, bufferedErr *bufio.Writer
}
func (o BufferedOutputs) Flush() error {
	outs := []io.Writer{o.bufferedLog, o.bufferedOut, o.bufferedErr}
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

func New(log, out, err io.Writer) BasicOutputs {
	t := time.Time{}
	aLog := ActivityWriter{log, t}
	aOut := ActivityWriter{out, t}
	aErr := ActivityWriter{err, t}
	return BasicOutputs{&aLog, &aOut, &aErr}
}

func NewStandardOutputs() Outputs {
	return New(os.Stdout, os.Stdout, os.Stderr)
}

func NewBufferedOutputs(outputs Outputs) Outputs {
	log := bufio.NewWriter(outputs.Log())
	out := bufio.NewWriter(outputs.Out())
	err := bufio.NewWriter(outputs.Err())
	basic := New(log, out, err)
	buffered := BufferedOutputs{basic, log, out, err}
	return buffered
}

// TODO ?
//func NewFileOutputs(outFilePath, errFilePath string) Outputs {
//	if outFilePath == errFilePath || errFilePath == "" {
//		// Same file for both outputs
//
//	}
//}
