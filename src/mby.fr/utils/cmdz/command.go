package cmdz

import (
	//"bufio"

	"bytes"
	"log"

	//"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"

	"mby.fr/utils/inout"
	"mby.fr/utils/promise"
	"mby.fr/utils/stringz"
)

type failure struct {
	Rc       int
	reporter Reporter
}

func (f failure) Error() (msg string) {
	if f.reporter != nil {
		stderrSummary := stringz.SummaryRatio(f.reporter.ReportError(), 128, 0.2)
		msg = fmt.Sprintf("Failing with ResultCode: %d executing: [%s] ! stderr: %s", f.Rc, f.reporter.String(), stderrSummary)
	}
	return
}

type config struct {
	retries        int
	retryDelayInMs int
	timeout        int
	stdin          io.Reader
	stdout         io.Writer
	stderr         io.Writer
	combinedOuts   bool
	errorOnFailure bool

	feeder      *cmdz
	pipedInput  bool
	pipedOutput bool
	pipeFail    bool
}

// Merge lower priority config into higher priority
func mergeConfigs(higher, lower *config) *config {
	merged := *higher
	if lower == nil {
		return &merged
	}
	if merged.retries == 0 {
		merged.retries = lower.retries
	}
	if merged.retryDelayInMs == 0 {
		merged.retryDelayInMs = lower.retryDelayInMs
	}
	if merged.timeout == 0 {
		merged.timeout = lower.timeout
	}
	if merged.stdout == nil {
		merged.stdout = lower.stdout
	}
	if merged.stderr == nil {
		merged.stderr = lower.stderr
	}
	if !merged.errorOnFailure {
		merged.errorOnFailure = lower.errorOnFailure
	}
	return &merged
}

type cmdz struct {
	*exec.Cmd
	config

	cmdCheckpoint  exec.Cmd
	fallbackConfig *config

	stdinRecord  inout.RecordingReader
	stdoutRecord inout.RecordingWriter
	stderrRecord inout.RecordingWriter

	inProcesser  inout.ProcessingReader
	outProcesser inout.ProcessingWriter
	errProcesser inout.ProcessingWriter

	resultsCodes []int
	// FIXME: replace ResultsCodes by Executions
	Executions []*exec.Cmd
}

func (e *cmdz) getConfig() config {
	return e.config
}

// ----- InOuter methods -----
func (e *cmdz) Stdin() io.Reader {
	return e.Cmd.Stdin
}
func (e *cmdz) Stdout() io.Writer {
	return e.Cmd.Stdout
}
func (e *cmdz) Stderr() io.Writer {
	return e.Cmd.Stdout
}

func (e *cmdz) SetInput(stdin io.Reader) Executer {
	if e.pipedInput {
		log.Fatal("Input is piped cannot change it !")
	}
	e.setupStdin(stdin)
	return e
}

func (e *cmdz) SetStdout(stdout io.Writer) Executer {
	if e.pipedOutput {
		log.Fatal("Output is piped cannot change it !")
	}
	e.setupStdout(stdout)
	return e
}

func (e *cmdz) SetStderr(stderr io.Writer) Executer {
	e.setupStderr(stderr)
	return e
}

func (e *cmdz) SetOutputs(stdout, stderr io.Writer) Executer {
	e.SetStdout(stdout)
	e.SetStderr(stderr)
	return e
}

func (e *cmdz) initProcessers() {
	if e.inProcesser == nil {
		e.inProcesser = inout.NewProcessingStreamReader(nil)
	}
	if e.outProcesser == nil {
		e.outProcesser = inout.NewProcessingStreamWriter(nil)
	}
	if e.errProcesser == nil {
		e.errProcesser = inout.NewProcessingStreamWriter(nil)
	}
}

func (e *cmdz) setupStdin(stdin io.Reader) {
	if stdin == nil {
		stdin = bytes.NewReader(nil)
	}
	// Save supplied stdin in config
	// Decorate reader: ProcessingReader => RecordingReader => stdin
	e.initProcessers()
	e.config.stdin = stdin
	e.inProcesser.Nest(stdin)
	e.stdinRecord.Nested = e.inProcesser
	e.Cmd.Stdin = &e.stdinRecord
}

func (e *cmdz) setupStdout(stdout io.Writer) {
	if stdout == nil {
		stdout = &bytes.Buffer{}
	}
	// Save supplied stdout in config
	// Decorate writer: ProcessingWriter => RecordingWriter => stdout
	e.initProcessers()
	e.config.stdout = stdout
	e.stdoutRecord.Nested = stdout
	e.outProcesser.Nest(&e.stdoutRecord)
	e.Cmd.Stdout = e.outProcesser
	if e.config.combinedOuts {
		e.Cmd.Stderr = e.outProcesser
	}
}

func (e *cmdz) setupStderr(stderr io.Writer) {
	if e.config.combinedOuts {
		return
	}
	if stderr == nil {
		stderr = &strings.Builder{}
	}
	// Save supplied stderr in config
	// Decorate writer: ProcessingWriter => RecordingWriter => stderr
	e.initProcessers()
	e.config.stderr = stderr
	e.stderrRecord.Nested = stderr
	e.errProcesser.Nest(&e.stderrRecord)
	e.Cmd.Stderr = e.errProcesser
}

func (e *cmdz) ProcessIn(pcrs ...IOProcesser) *cmdz {
	e.initProcessers()
	e.inProcesser.Add(pcrs...)
	return e
}

func (e *cmdz) ProcessOut(pcrs ...IOProcesser) *cmdz {
	e.initProcessers()
	e.outProcesser.Add(pcrs...)
	return e
}

func (e *cmdz) ProcessErr(pcrs ...IOProcesser) *cmdz {
	e.initProcessers()
	e.errProcesser.Add(pcrs...)
	return e
}

// ----- Configurer methods -----
func (e *cmdz) ErrorOnFailure(enable bool) Executer {
	e.config.errorOnFailure = enable
	return e
}

func (e *cmdz) Retries(count, delayInMs int) Executer {
	e.config.retries = count
	e.config.retryDelayInMs = delayInMs
	return e
}

func (e *cmdz) Timeout(delayInMs int) Executer {
	e.config.timeout = delayInMs
	return e
}

func (e *cmdz) CombinedOutputs() Executer {
	e.config.combinedOuts = true
	return e
}

// ----- Recorder methods -----
func (e *cmdz) StdinRecord() string {
	return e.stdinRecord.String()
}

func (e *cmdz) StdoutRecord() string {
	return e.stdoutRecord.String()
}

func (e *cmdz) StderrRecord() string {
	return e.stderrRecord.String()
}

// ----- Reporter methods -----
func (e cmdz) String() (t string) {
	t = strings.Join(e.Args, " ")
	return
}

func (e cmdz) ReportError() string {
	execCmdSummary := e.String()
	attempts := len(e.resultsCodes)
	status := e.resultsCodes[attempts-1]
	stderr := e.stderrRecord.String()
	errorMessage := fmt.Sprintf("Exec failed after %d attempt(s): [%s] !\nRC=%d ERR> %s", attempts, execCmdSummary, status, strings.TrimSpace(stderr))
	return errorMessage
}

// ----- Runner methods -----
func (e *cmdz) reset() {
	e.stdinRecord.Reset()
	e.stdoutRecord.Reset()
	e.stderrRecord.Reset()
	e.resultsCodes = nil
	e.Executions = nil
	e.rollback()
}

func (e *cmdz) fallback(cfg *config) {
	e.fallbackConfig = cfg
}

func (e *cmdz) BlockRun() (rc int, err error) {
	f := e.feeder
	var originalStdout io.Writer
	var originalStdin io.Reader
	if f != nil {
		originalStdout = f.config.stdout
		originalStdin = f.config.stdin
		b := bytes.Buffer{}

		// Replace configured stdin / stdout temporarilly
		f.pipedOutput = true
		f.config.stdout = &b

		e.pipedInput = true
		e.config.stdin = &b

		frc, ferr := e.feeder.BlockRun()
		if _, ok := ferr.(failure); ferr != nil && !ok {
			// If feeder error is not a failure{} return it immediately
			return frc, ferr
		}

		if e.feeder.pipeFail && frc > 0 {
			// if pipefail enabled and rc > 0 fail shortly
			return frc, ferr
		}
	}

	e.reset()
	config := mergeConfigs(&e.config, e.fallbackConfig)
	e.setupStdin(config.stdin)
	e.setupStdout(config.stdout)
	e.setupStderr(config.stderr)
	e.checkpoint()
	rc = -1
	for i := 0; i <= config.retries && rc != 0; i++ {

		if i > 0 {
			// Wait between retries
			time.Sleep(time.Duration(config.retryDelayInMs) * time.Millisecond)
		}
		if commandMock == nil {
			err = e.Start()
			if err != nil {
				return
			}
			err = e.Wait()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					rc = exitErr.ProcessState.ExitCode()
					err = nil
				} else {
					return -1, err
				}
			} else {
				rc = e.ProcessState.ExitCode()
			}
		} else {
			// Replace execution by mocking function
			rc = commandMock.Mock(e.Cmd)
		}
		e.resultsCodes = append(e.resultsCodes, rc)
		e.Executions = append(e.Executions, e.Cmd)
		e.rollback()
	}
	if e.errorOnFailure && rc > 0 {
		err = failure{rc, e}
		rc = -1
	}
	if f != nil {
		f.config.stdout = originalStdout
		e.config.stdin = originalStdin
	}
	return
}

func (e *cmdz) AsyncRun() *execPromise {
	p := promise.New(func(resolve func(int), reject func(error)) {
		rc, err := e.BlockRun()
		if err != nil {
			reject(err)
		}
		if e.errorOnFailure && rc > 0 {
			err = failure{rc, e}
			reject(err)
		}
		resolve(rc)
	})
	return p
}

func (e *cmdz) ResultCodes() []int {
	return e.resultsCodes
}

// ----- Piper methods -----
func (e *cmdz) Pipe(c *cmdz) *cmdz {
	c.feeder = e
	return c
}

func (e *cmdz) PipeFail(c *cmdz) *cmdz {
	e.pipeFail = true
	return e.Pipe(c)
}

func (e *cmdz) AddEnv(key, value string) *cmdz {
	entry := fmt.Sprintf("%s=%s", key, value)
	e.Cmd.Env = append(e.Env, entry)
	e.checkpoint()
	return e
}

func (e *cmdz) AddArgs(args ...string) *cmdz {
	e.Cmd.Args = append(e.Args, args...)
	e.checkpoint()
	return e
}

func (e *cmdz) checkpoint() {
	e.cmdCheckpoint = *e.Cmd
}

func (e *cmdz) rollback() {
	// Clone checkpoint
	newClone := e.cmdCheckpoint
	e.Cmd = &newClone
}