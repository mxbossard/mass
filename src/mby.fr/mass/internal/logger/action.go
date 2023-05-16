package logger

import (
	"fmt"
	"sync"

	"mby.fr/mass/internal/output"
	"mby.fr/utils/ansi"
	"mby.fr/utils/format"
	"mby.fr/utils/inout"
	"mby.fr/utils/logz"
)

const (
	actionPadding = 30
)

var (
	outAnsiColors []string = []string{ansi.Reset, ansi.HiGreen, ansi.HiBlue, ansi.HiCyan, ansi.HiWhite, ansi.Green, ansi.Blue, ansi.Cyan, ansi.White}
	errAnsiColors []string = []string{ansi.HiRed, ansi.HiYellow, ansi.HiPurple, ansi.Red, ansi.Yellow, ansi.Purple}
)

type ActionLogger struct {
	logz.Logger
	output.Outputs
}

func (l ActionLogger) Start() {
}

func (l ActionLogger) End() {
}

func (l ActionLogger) Progress() {
}

func NewAction(outs output.Outputs, action, subject string, filterLevel int) ActionLogger {
	loggerName := fmt.Sprintf("%s(%s)", action, subject)

	// Decorate outputs
	outColorFormatter := inout.AnsiFormatter{getOutAnsiColor()}
	errColorFormatter := inout.AnsiFormatter{getErrAnsiColor()}
	outPrefixedFormatter := inout.PrefixFormatter{Prefix: "out>", RightPad: 5}
	errPrefixedFormatter := inout.PrefixFormatter{Prefix: "err>", RightPad: 5}

	loggerPrefixedFormatter := inout.LineFormatter{func(line string) string {
		prefix := fmt.Sprintf("[%s] ", format.PadRight(loggerName, actionPadding))
		return prefix + line
	}}

	//log := outs.Log()
	log := outs.Out()
	log = inout.NewFormattingWriter(log, outColorFormatter)
	out := outs.Out()
	out = inout.NewFormattingWriter(out, outColorFormatter)
	out = inout.NewFormattingWriter(out, loggerPrefixedFormatter)
	out = inout.NewFormattingWriter(out, outPrefixedFormatter)
	err := outs.Err()
	err = inout.NewFormattingWriter(err, errColorFormatter)
	err = inout.NewFormattingWriter(err, loggerPrefixedFormatter)
	err = inout.NewFormattingWriter(err, errPrefixedFormatter)
	decoratedOuts := output.New(log, out, err)

	logger := logz.New(log, loggerName, actionPadding, true, false, filterLevel)
	al := ActionLogger{logger, decoratedOuts}
	return al
}

var outAnsiColorCounter = 0
var outAnsiColorMutex = sync.Mutex{}
var errAnsiColorCounter = 0
var errAnsiColorMutex = sync.Mutex{}

func getOutAnsiColor() string {
	outAnsiColorMutex.Lock()
	defer outAnsiColorMutex.Unlock()
	ansiColor := outAnsiColors[outAnsiColorCounter%len(outAnsiColors)]
	outAnsiColorCounter++
	return ansiColor
}

func getErrAnsiColor() string {
	errAnsiColorMutex.Lock()
	defer errAnsiColorMutex.Unlock()
	ansiColor := errAnsiColors[errAnsiColorCounter%len(errAnsiColors)]
	errAnsiColorCounter++
	return ansiColor
}
