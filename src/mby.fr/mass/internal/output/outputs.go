package output

import (
	"io"
	"os"
	"bufio"
	"time"
	//"fmt"

	"mby.fr/utils/datetime"
	"mby.fr/utils/inout"
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

type ActivityOutputs struct {
	log, out, err *inout.ActivityWriter
}

func (o ActivityOutputs) Flush() error {
	outs := []io.Writer{o.log.Nested, o.out.Nested, o.err.Nested}
	for _, out := range outs {
		f, ok := out.(Flusher)
		if ok {
			err := f.Flush()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o ActivityOutputs) LastWriteTime() time.Time {
	lastTime := datetime.Max(o.log.Activity, o.out.Activity, o.err.Activity)
	return lastTime
}

func (o ActivityOutputs) Log() io.Writer {
	return o.log
}

func (o ActivityOutputs) Out() io.Writer {
	return o.out
}

func (o ActivityOutputs) Err() io.Writer {
	return o.err
}

type AwareOutputs struct {
	ActivityOutputs
	underlying *Outputs
}

func (o AwareOutputs) Flush() error {
	return (*o.underlying).Flush()
}

func New(log, out, err io.Writer) ActivityOutputs {
	t := time.Time{}
	aLog := inout.ActivityWriter{log, t}
	aOut := inout.ActivityWriter{out, t}
	aErr := inout.ActivityWriter{err, t}
	return ActivityOutputs{&aLog, &aOut, &aErr}
}

func NewStandardOutputs() Outputs {
	return New(os.Stdout, os.Stdout, os.Stderr)
}

func NewBufferedOutputs(outputs Outputs) Outputs {
	log := bufio.NewWriter(outputs.Log())
	out := bufio.NewWriter(outputs.Out())
	err := bufio.NewWriter(outputs.Err())
	buffered := New(log, out, err)
	return buffered
}

// TODO ?
//func NewFileOutputs(outFilePath, errFilePath string) Outputs {
//	if outFilePath == errFilePath || errFilePath == "" {
//		// Same file for both outputs
//
//	}
//}
