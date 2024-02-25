package display

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"mby.fr/cmdtest/model"
	"mby.fr/cmdtest/utils"
	"mby.fr/utils/ansi"
	"mby.fr/utils/cmdz"
	"mby.fr/utils/inout"
	"mby.fr/utils/printz"
)

var (
	messageColor = ansi.HiPurple
	testColor    = ansi.HiCyan
	successColor = ansi.BoldGreen
	failureColor = ansi.BoldRed
	reportColor  = ansi.Yellow
	warningColor = ansi.BoldHiYellow
	errorColor   = ansi.Red
)

type Displayer interface {
	Global(model.Context)
	Suite(model.Context)
	TestTitle(model.Context)
	TestOutcome(ctx model.Context)
	AssertionResult(model.AssertionResult)
	ReportSuite(model.Context)
	ReportAll(model.Context)
	Stdout(string)
	Stderr(string)
	Error(error)
	Flush() error
}

type BasicDisplay struct {
	printer            printz.Printer
	clearAnsiFormatter inout.Formatter
	outFormatter       inout.Formatter
	errFormatter       inout.Formatter
}

func (d BasicDisplay) Global(ctx model.Context) {
	defer d.Flush()
	// Do nothing ?
	if ctx.Config.Verbose.Get() >= model.BETTER_ASSERTION_REPORT {
		d.printer.ColoredErrf(messageColor, "## New config (token: %s)\n", ctx.Token)
	}
}

func (d BasicDisplay) Suite(ctx model.Context) {
	defer d.Flush()
	if ctx.Config.Verbose.Get() >= model.BETTER_ASSERTION_REPORT {
		d.printer.ColoredErrf(messageColor, "## New test suite: [%s] (token: %s)\n", ctx.TestQualifiedName(), ctx.Token)
	}
}

func (d BasicDisplay) TestTitle(ctx model.Context, seq int) {
	// FIXME: get seq from ctx
	defer d.Flush()
	timecode := int(time.Since(ctx.StartTime).Milliseconds())
	qulifiedName := ctx.TestName
	if ctx.TestSuite != "" {
		qulifiedName = fmt.Sprintf("[%s]/%s", ctx.TestSuite, ctx.TestName)
	}

	if ctx.Ignore != nil && *ctx.Ignore {
		if ctx.Silent == nil || !*ctx.Silent {
			d.printer.ColoredErrf(warningColor, "[%05d] Test: %s #%02d... ", timecode, qulifiedName, seq)
		}
		return
	}

	testTitle := fmt.Sprintf("[%05d] Test %s #%02d", timecode, qulifiedName, seq)
	if ctx.Silent == nil || !*ctx.Silent {
		d.printer.ColoredErrf(testColor, "%s... ", testTitle)
		if *ctx.KeepStdout || *ctx.KeepStderr {
			// NewLine because we expect cmd outputs
			d.printer.Errf("\n")
		}
	}
}

func (d BasicDisplay) TestOutcome(ctx model.Context, seq int, outcome model.Outcome, cmd cmdz.Executer, testDuration time.Duration, err error) {
	// FIXME get outcome from ctx
	// FIXME get cmd, duration and error from outcome
	defer d.Flush()
	switch outcome {
	case model.PASSED:
		if ctx.Silent == nil || !*ctx.Silent {
			d.printer.ColoredErrf(successColor, "PASSED")
			d.printer.Errf(" (in %s)\n", testDuration)
		}
	case model.FAILED:
		if ctx.Silent != nil && *ctx.Silent {
			*ctx.Silent = false
			d.TestTitle(ctx, 0)
		}
		if err == nil {
			//IncrementSeq(tmpDir, FailureSequenceFilename)
			d.printer.ColoredErrf(failureColor, "FAILED")
			d.printer.Errf(" (in %s)\n", testDuration)
		} else {
			if errors.Is(err, context.DeadlineExceeded) {
				// Swallow error
				err = nil
				//IncrementSeq(tmpDir, FailureSequenceFilename)
				d.printer.ColoredErrf(failureColor, "FAILED")
				d.printer.Errf(" (timed out after %s)\n", ctx.Timeout)
			}
		}
		d.printer.Errf("Failure calling cmd: <|%s|>\n", cmd)
	case model.ERRORED:
		d.printer.ColoredErrf(warningColor, "ERROR")
		d.printer.Errf(" (not executed)\n")
		d.printer.Errf("Failure calling cmd: <|%s|>\n", cmd)
	case model.IGNORED:
		if ctx.Silent == nil || !*ctx.Silent {
			d.printer.ColoredErrf(warningColor, "IGNORED")
			d.printer.Err("\n")
		}
	default:
	}
	if err != nil {
		d.printer.ColoredErrf(model.ErrorColor, "%s\n", err)
	}

	if (ctx.Silent == nil || !*ctx.Silent) && (*ctx.KeepStdout || *ctx.KeepStderr) {
		// NewLine in printer to print test result in a new line
		d.printer.Errf("        ")
		d.printer.Flush()
	}
}

func (d BasicDisplay) AssertionResult(cmd cmdz.Executer, result model.AssertionResult) {
	defer d.Flush()
	//log.Printf("failedResult: %v\n", result)
	assertPrefix := result.Assertion.Prefix
	assertName := result.Assertion.Name
	assertOp := result.Assertion.Op
	expected := result.Assertion.Expected
	got := result.Value

	if result.Message != "" {
		d.printer.ColoredErrf(errorColor, result.Message+"\n")
	}

	if assertName == "success" || assertName == "fail" {
		d.printer.Errf("Expected %s%s\n", assertPrefix, assertName)
		d.Stdout(cmd.StdoutRecord())
		d.Stderr(cmd.StderrRecord())
		/*
			if cmd.StderrRecord() != "" {
				d.printer.Errf("sdterr> %s\n", cmd.StderrRecord())
			}
		*/
		return
	} else if assertName == "cmd" {
		d.printer.Errf("Expected %s%s=%s to succeed\n", assertPrefix, assertName, expected)
		return
	} else if assertName == "exists" {
		d.printer.Errf("Expected file %s%s=%s file to exists\n", assertPrefix, assertName, expected)
		return
	}

	if expected != got {

		if s, ok := got.(string); ok {
			got = strings.ReplaceAll(s, "\n", "\\n")
			/*
				const sliceSize = 16
				minStrLen := min(len(s), len(expected))
				for k := range minStrLen / sliceSize {
					left := expected[k*sliceSize : min(len(expected), (k+1)*sliceSize-1)]
					right := s[k*sliceSize : min(len(s), (k+1)*sliceSize-1)]
					if left == right {
						continue
					} else {
						shortenExpected := ""
						if k > 0 {
							shortenExpected += "[...]"
						}
						shortenExpected += left
						if k*minStrLen < len(expected) {
							shortenExpected += "[...]"
						}
						shortenGot := ""
						if k > 0 {
							shortenGot += "[...]"
						}
						shortenGot += right
						if k*minStrLen < len(s) {
							shortenGot += "[...]"
						}
						expected = shortenExpected
						s = shortenGot
					}
				}
					expected = d.clearAnsiFormatter.Format(expected)
					got = d.clearAnsiFormatter.Format(s)
			*/
		}
		if assertOp == "=" {
			d.printer.Errf("Expected %s%s to be: [%s] but got: [%v]\n", assertPrefix, assertName, expected, got)
		} else if assertOp == ":" {
			d.printer.Errf("Expected %s%s to contains: [%s] but got: [%v]\n", assertPrefix, assertName, expected, got)
		} else if assertOp == "!:" {
			d.printer.Errf("Expected %s%s not to contains: [%s] but got: [%v]\n", assertPrefix, assertName, expected, got)
		} else if assertOp == "~" {
			d.printer.Errf("Expected %s%s to match: [%s] but got: [%v]\n", assertPrefix, assertName, expected, got)
		} else if assertOp == "!~" {
			d.printer.Errf("Expected %s%s not to match: [%s] but got: [%v]\n", assertPrefix, assertName, expected, got)
		}
	} else {
		d.printer.Errf("assertion %s%s%s%s failed\n", assertPrefix, assertName, assertOp, expected)
	}
}

func (d BasicDisplay) ReportSuite(ctx model.Context, tmpDir string, failedReports []string) {
	defer d.Flush()
	// FIXME: get tmpDir and failed reports from ctx
	testCount := utils.ReadSeq(tmpDir, model.TestSequenceFilename)       // TODO: put in model.Context
	ignoredCount := utils.ReadSeq(tmpDir, model.IgnoredSequenceFilename) // TODO: put in model.Context
	failedCount := utils.ReadSeq(tmpDir, model.FailedSequenceFilename)   // TODO: put in model.Context
	errorCount := utils.ReadSeq(tmpDir, model.ErroredSequenceFilename)   // TODO: put in model.Context
	passedCount := testCount - failedCount - ignoredCount

	if ctx.Silent == nil || !*ctx.Silent {
		d.printer.ColoredErrf(messageColor, "Reporting [%s] test suite (%s) ...\n", ctx.TestSuite, tmpDir)
	}

	ignoredMessage := ""
	if ignoredCount > 0 {
		ignoredMessage = fmt.Sprintf(" (%d ignored)", ignoredCount)
	}
	duration := ctx.LastTestTime.Sub(ctx.StartTime)
	fmtDuration := NormalizeDurationInSec(duration)
	if failedCount == 0 && errorCount == 0 {
		d.printer.ColoredErrf(successColor, "Successfuly ran [%s] test suite (%d tests in %s)", ctx.TestSuite, passedCount, fmtDuration)
		d.printer.ColoredErrf(warningColor, "%s", ignoredMessage)
		d.printer.Errf("\n")
	} else {
		d.printer.ColoredErrf(failureColor, "Failures in [%s] test suite (%d success, %d failures, %d errors on %d tests in %s)", ctx.TestSuite, passedCount, failedCount, errorCount, testCount, fmtDuration)
		d.printer.ColoredErrf(warningColor, "%s", ignoredMessage)
		d.printer.Errf("\n")
		for _, report := range failedReports {
			d.printer.ColoredErrf(reportColor, "%s\n", report)
		}
	}
}

func (d BasicDisplay) ReportAllFooter(testSuitesCtx model.Context) {
	defer d.Flush()

	globalStartTime := time.Now() // FIXME
	globalDuration := model.NormalizeDurationInSec(time.Since(globalStartTime))
	d.printer.ColoredErrf(reportColor, "Global duration time: %s\n", globalDuration)
}

func (d BasicDisplay) Stdout(s string) {
	if s != "" {
		d.printer.Err(d.outFormatter.Format(s))
	}
}

func (d BasicDisplay) Stderr(s string) {
	if s != "" {
		d.printer.Err(d.errFormatter.Format(s))
	}
}

func (d BasicDisplay) Error(err error) {

}

func (d BasicDisplay) Flush() error {
	return d.printer.Flush()
}

func New() BasicDisplay {
	return BasicDisplay{
		printer:            printz.NewStandard(),
		clearAnsiFormatter: inout.AnsiFormatter{AnsiFormat: ansi.Reset},
		outFormatter:       inout.PrefixFormatter{Prefix: "out> "},
		errFormatter:       inout.PrefixFormatter{Prefix: "err> "},
	}
}

func NormalizeDurationInSec(d time.Duration) (duration string) {
	duration = fmt.Sprintf("%.3f s", float32(d.Milliseconds())/1000)
	return
}
