package display

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"mby.fr/cmdtest/model"
	"mby.fr/cmdtest/service"
	"mby.fr/utils/ansi"
	"mby.fr/utils/cmdz"
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
	AssertionResult(service.AssertionResult)
	ReportAll(model.Context)
	ReportSuite(model.Context)
	Stdout(string)
	Stderr(string)
	Error(error)
	Flush() error
}

type BasicDisplay struct {
	printer printz.Printer
}

func (d BasicDisplay) Global(ctx service.Context) {
	defer d.Flush()
	// Do nothing ?
	if ctx.Silent == nil || !*ctx.Silent {
		d.printer.ColoredErrf(messageColor, "Initialized new config [%s].\n", ctx.TestSuite)
	}
}

func (d BasicDisplay) Suite(ctx service.Context) {
	defer d.Flush()
	if ctx.Silent == nil || !*ctx.Silent {
		var tokenMsg = ""
		if ctx.Token != "" {
			tokenMsg = fmt.Sprintf(" (token: %s)", ctx.Token)
		}
		d.printer.ColoredErrf(messageColor, "Initialized new [%s] test suite%s.\n", ctx.TestSuite, tokenMsg)
	}
}

func (d BasicDisplay) TestTitle(ctx service.Context, seq int) {
	// FIXME: get seq from ctx
	defer d.Flush()
	timecode := int(time.Since(ctx.StartTime).Milliseconds())
	qulifiedName := ctx.TestName
	if ctx.TestSuite != "" {
		qulifiedName = fmt.Sprintf("[%s]/%s", ctx.TestSuite, ctx.TestName)
	}

	if ctx.Ignore != nil && *ctx.Ignore {
		if ctx.Silent == nil || !*ctx.Silent {
			d.printer.ColoredErrf(warningColor, "[%05d] Ignored test: %s\n", timecode, qulifiedName)
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

func (d BasicDisplay) TestOutcome(ctx service.Context, outcome string, cmd cmdz.Executer, testDuration, err error) {
	// FIXME get outcome from ctx
	// FIXME get cmd, duration and error from outcome
	defer d.Flush()
	switch outcome {
	case "PASSED":
		if ctx.Silent == nil || !*ctx.Silent {
			d.printer.ColoredErrf(successColor, "PASSED")
			d.printer.Errf(" (in %s)\n", testDuration)
		}
	case "FAILED":
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
	case "ERRORED":
		d.printer.ColoredErrf(warningColor, "ERROR")
		d.printer.Errf(" (not executed)\n")
		d.printer.Errf("Failure calling cmd: <|%s|>\n", cmd)
	case "IGNORED":
	default:
	}
}

func (d BasicDisplay) AssertionResult(cmd cmdz.Executer, result service.AssertionResult) {
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
		if cmd.StderrRecord() != "" {
			d.printer.Errf("sdterr> %s\n", cmd.StderrRecord())
		}
		return
	} else if assertName == "cmd" {
		d.printer.Errf("Expected %s%s=%s to succeed\n", assertPrefix, assertName, expected)
		return
	} else if assertName == "exists" {
		d.printer.Errf("Expected file %s%s=%s file to exists\n", assertPrefix, assertName, expected)
		return
	}

	if expected != got {
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

func (d BasicDisplay) ReportAll(ctx service.Context) {
	defer d.Flush()
	// TODO: list test suite and iterate over it

	globalStartTime := time.Now() // FIXME
	globalDuration := service.NormalizeDurationInSec(time.Since(globalStartTime))
	d.printer.ColoredErrf(reportColor, "Global duration time: %s\n", globalDuration)
}

func (d BasicDisplay) ReportSuite(ctx service.Context) {
	tmpDir, err := service.TestsuiteDirectoryPath(ctx.TestSuite, ctx.Token)
	if err != nil {
		log.Fatal(err)
	}
	testCount := service.ReadSeq(tmpDir, service.TestSequenceFilename)       // TODO: put in model.Context
	ignoredCount := service.ReadSeq(tmpDir, service.IgnoredSequenceFilename) // TODO: put in model.Context
	failureCount := service.ReadSeq(tmpDir, service.FailureSequenceFilename) // TODO: put in model.Context
	errorCount := service.ReadSeq(tmpDir, service.ErrorSequenceFilename)     // TODO: put in model.Context
	failedReports, err := service.FailureReports(ctx.TestSuite, ctx.Token)   // TODO: put in model.Context
	if err != nil {
		log.Fatal(err)
	}
	//failedCount := len(failedReports)
	failedCount := failureCount

	if ctx.Silent == nil || !*ctx.Silent {
		d.printer.ColoredErrf(messageColor, "Reporting [%s] test suite (%s) ...\n", ctx.TestSuite, tmpDir)
	}

	ignoredMessage := ""
	if ignoredCount > 0 {
		ignoredMessage = fmt.Sprintf(" (%d ignored)", ignoredCount)
	}
	duration := ctx.LastTestTime.Sub(ctx.StartTime)
	fmtDuration := service.NormalizeDurationInSec(duration)
	if failedCount == 0 && errorCount == 0 {
		d.printer.ColoredErrf(successColor, "Successfuly ran [%s] test suite (%d tests in %s)", ctx.TestSuite, testCount, fmtDuration)
		d.printer.ColoredErrf(warningColor, "%s", ignoredMessage)
		d.printer.Errf("\n")
	} else {
		successCount := testCount - failedCount
		d.printer.ColoredErrf(failureColor, "Failures in [%s] test suite (%d success, %d failures, %d errors on %d tests in %s)", ctx.TestSuite, successCount, failedCount, errorCount, testCount, fmtDuration)
		d.printer.ColoredErrf(warningColor, "%s", ignoredMessage)
		d.printer.Errf("\n")
		for _, report := range failedReports {
			d.printer.ColoredErrf(reportColor, "%s\n", report)
		}
	}
}

func (d BasicDisplay) Stdout(s string) {

}

func (d BasicDisplay) Stderr(s string) {

}

func (d BasicDisplay) Error(err error) {

}

func (d BasicDisplay) Flush() error {
	return d.printer.Flush()
}

func New() BasicDisplay {
	return BasicDisplay{printer: printz.NewStandard()}
}
