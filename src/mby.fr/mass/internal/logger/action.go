package logger

import (
	"fmt"
	"sync"

	"mby.fr/mass/internal/output"
	"mby.fr/utils/logger"
	"mby.fr/utils/ansi"
	"mby.fr/utils/inout"
)

var (
	outAnsiColors []string = []string{ansi.Reset, ansi.HiGreen, ansi.HiBlue, ansi.HiCyan, ansi.HiWhite, ansi.Green, ansi.Blue, ansi.Cyan, ansi.White}
	errAnsiColors []string = []string{ansi.HiRed, ansi.HiYellow, ansi.HiPurple, ansi.Red, ansi.Yellow, ansi.Purple}
)


type ActionLogger struct {
	logger.Logger
	output.Outputs
}

func (l ActionLogger) Start() {
}

func (l ActionLogger) End() {
}

func (l ActionLogger) Progress() {
}

func NewAction(outs output.Outputs, action, subject string) ActionLogger {
	loggerName := fmt.Sprintf("%s(%s)", action, subject)

	// Decorate outputs
	outColorFormatter := inout.AnsiFormatter{getOutAnsiColor()}
	errColorFormatter := inout.AnsiFormatter{getErrAnsiColor()}
	prefixedFormatter := inout.PrefixFormatter{Prefix: loggerName + " ", RightPad: 25}
	outPrefixedFormatter := inout.PrefixFormatter{Prefix: ">O "}
	errPrefixedFormatter := inout.PrefixFormatter{Prefix: ">E "}

	log := outs.Log()
	log = inout.NewFormattingWriter(log, outColorFormatter)
	out := outs.Out()
	out = inout.NewFormattingWriter(out, outColorFormatter)
	out = inout.NewFormattingWriter(out, outPrefixedFormatter)
	out = inout.NewFormattingWriter(out, prefixedFormatter)
	err := outs.Err()
	err = inout.NewFormattingWriter(err, errColorFormatter)
	err = inout.NewFormattingWriter(err, errPrefixedFormatter)
	err = inout.NewFormattingWriter(err, prefixedFormatter)
	decoratedOuts := output.New(log, out, err)

	logger := logger.New(log, loggerName, true, true)
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
        ansiColor := outAnsiColors[outAnsiColorCounter % len(outAnsiColors)]
        outAnsiColorCounter ++
        return ansiColor
}

func getErrAnsiColor() string {
	errAnsiColorMutex.Lock()
	defer errAnsiColorMutex.Unlock()
        ansiColor := errAnsiColors[errAnsiColorCounter % len(errAnsiColors)]
        errAnsiColorCounter ++
        return ansiColor
}
