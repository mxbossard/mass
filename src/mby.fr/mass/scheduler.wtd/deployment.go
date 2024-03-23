package scheduler

import (
	"io"
)

type Event struct {
	Timestamp   int64
	Name        string
	Description string
}

type Eventer interface {
	Events() []*Event
}

type Executioner interface {
	WaitCompletion() error
	ExitCode() (int, error)
	StdOut() io.Reader
	StdErr() io.Reader
}

type Execution struct {
	exitCode int
	stdOut   io.Reader
	stdErr   io.Reader
}

type Container interface {
	Eventer
	Run(args []string) (Executioner, error)
	Stop() error
}

type Probe interface {
	Eventer
	Probe() (bool, error)
}

type Deployment struct {
	Eventer
}

func (d Deployment) Start() (err error) {
	return
}

func (d Deployment) Events() (events []*Event) {
	return
}
