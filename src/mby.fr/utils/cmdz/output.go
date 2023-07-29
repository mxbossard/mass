package cmdz

import (
	"fmt"
)

type (
	basicOutput struct {
		Executer
		inProcessers  []InProcesser
		outProcessers []OutProcesser
	}
)

// For simplicity processers only apply to Output() and OutputString()
// FIXME: processers should be embedded in decorating Writer in front of stdout
func (e *basicOutput) OutProcess(processers ...OutProcesser) Outputer {
	e.outProcessers = append(e.outProcessers, processers...)
	return e
}

func (e *basicOutput) OutStringProcess(strProcessers ...OutStringProcesser) Outputer {
	processers := make([]OutProcesser, len(strProcessers))
	for i, sp := range strProcessers {
		spf := sp // Need this to update func pointer
		processers[i] = func(rc int, stdout, stderr []byte) ([]byte, error) {
			res, err := spf(rc, string(stdout), string(stderr))
			return []byte(res), err
		}

	}
	e.OutProcess(processers...)
	return e
}

func (e *basicOutput) Output() ([]byte, error) {
	return nil, fmt.Errorf("basicOutput.Output() not implemented yet !")
}

func (e *basicOutput) OutputString() (string, error) {
	rc, err := e.Executer.ErrorOnFailure(true).BlockRun()
	if err != nil {
		return "", err
	}
	if rc > 0 {
		err = failure{rc, e}
		return "", err
	}
	stdout := []byte(e.StdoutRecord())
	stderr := []byte(e.StderrRecord())
	for _, p := range e.outProcessers {
		stdout, err = p(rc, stdout, stderr)
		if err != nil {
			return "", err
		}
	}
	return string(stdout), nil
}

func (e *basicOutput) CombinedOutput() ([]byte, error) {
	return nil, fmt.Errorf("basicOutput.CombinedOutput() not implemented yet !")
}

func (e *basicOutput) CombinedOutputString() (string, error) {
	return e.CombinedOutputs().OutputString()
}

// ----- Override Configurer methods -----
func (e *basicOutput) ErrorOnFailure(ok bool) Outputer {
	e.Executer = e.Executer.ErrorOnFailure(ok)
	return e
}

func (e *basicOutput) CombinedOutputs() Outputer {
	e.Executer = e.Executer.CombinedOutputs()
	return e
}

func (e *basicOutput) Retries(count, delayInMs int) Outputer {
	e.Executer = e.Executer.Retries(count, delayInMs)
	return e
}

func (e *basicOutput) Timeout(delayInMs int) Outputer {
	e.Executer = e.Executer.Timeout(delayInMs)
	return e
}

func Outputted(e Executer) Outputer {
	return &basicOutput{Executer: e}
}

func OutputtedCmd(binary string, args ...string) Outputer {
	return Outputted(Cmd(binary, args...))
}
