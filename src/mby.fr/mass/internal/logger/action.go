package logger

import (
	"fmt"
	"sync"

	"mby.fr/mass/internal/output"
	"mby.fr/utils/anzi"
	"mby.fr/utils/formatz"
	"mby.fr/utils/inoutz"
	"mby.fr/utils/logz_toDel"
	"mby.fr/utils/ztring"
)

const (
	actionPadding = 40
)

var (
	outAnsiColors []anzi.Color = []anzi.Color{anzi.Reset, anzi.HiGreen, anzi.HiBlue, anzi.HiCyan, anzi.HiWhite, anzi.Green, anzi.Blue, anzi.Cyan, anzi.White}
	errAnsiColors []anzi.Color = []anzi.Color{anzi.HiRed, anzi.HiYellow, anzi.HiPurple, anzi.Red, anzi.Yellow, anzi.Purple}
)

type ActionLogger struct {
	logz_toDel.Logger
	output.Outputs
}

func (l ActionLogger) Start() {
}

func (l ActionLogger) End() {
}

func (l ActionLogger) Progress() {
}

func forgeLoggerName(action, subject string) (loggerName string) {
	loggerName = fmt.Sprintf("%s> %s", action, subject)
	return
}

func forgeActionPrefix(action, subject string) (actionPrefix string) {
	actionPrefix = forgeLoggerName(action, subject)
	if actionPadding < len(actionPrefix) {
		// loggerName too long to be displayed
		subjectParts, separators := ztring.SplitByRegexp(subject, "[ /,;:]")
		subjectMaxSize := actionPadding - len(forgeLoggerName(action, "")) - len(separators)
		subjectPartSize := subjectMaxSize / len(subjectParts)
		shortenedSubject := ""
		for k, sep := range separators {
			shortenedSubject += ztring.Left(subjectParts[k], subjectPartSize)
			shortenedSubject += sep
		}
		lastSubjectPart := subjectParts[len(subjectParts)-1]
		remainingSpace := subjectMaxSize + len(separators) - len(shortenedSubject)
		shortenedSubject += ztring.Left(lastSubjectPart, remainingSpace)
		actionPrefix = forgeLoggerName(action, shortenedSubject)
	}
	return
}

func NewAction(outs output.Outputs, action, subject string, filterLevel int) ActionLogger {
	loggerName := forgeLoggerName(action, subject)
	actionPrefix := forgeActionPrefix(action, subject)

	// Decorate outputs
	outColorFormatter := inoutz.AnsiFormatter{AnsiFormat: getOutAnsiColor()}
	errColorFormatter := inoutz.AnsiFormatter{AnsiFormat: getErrAnsiColor()}
	outPrefixedFormatter := inoutz.PrefixFormatter{Prefix: "out>", RightPad: 5}
	errPrefixedFormatter := inoutz.PrefixFormatter{Prefix: "err>", RightPad: 5}

	loggerPrefixedFormatter := inoutz.LineFormatter{Olf: func(line string) string {
		prefix := fmt.Sprintf("%s |", formatz.PadRight(actionPrefix, actionPadding))
		return prefix + line
	}}

	//log := outs.Log()
	log := outs.Out()
	log = inoutz.NewFormattingWriter(log, outColorFormatter)
	out := outs.Out()
	out = inoutz.NewFormattingWriter(out, outColorFormatter)
	out = inoutz.NewFormattingWriter(out, loggerPrefixedFormatter)
	out = inoutz.NewFormattingWriter(out, outPrefixedFormatter)
	err := outs.Err()
	err = inoutz.NewFormattingWriter(err, errColorFormatter)
	err = inoutz.NewFormattingWriter(err, loggerPrefixedFormatter)
	err = inoutz.NewFormattingWriter(err, errPrefixedFormatter)
	decoratedOuts := output.New(log, out, err)

	logger := logz_toDel.New(log, loggerName, actionPadding, true, false, filterLevel)
	al := ActionLogger{logger, decoratedOuts}
	return al
}

var outAnsiColorCounter = 0
var outAnsiColorMutex = sync.Mutex{}
var errAnsiColorCounter = 0
var errAnsiColorMutex = sync.Mutex{}

func getOutAnsiColor() anzi.Color {
	outAnsiColorMutex.Lock()
	defer outAnsiColorMutex.Unlock()
	ansiColor := outAnsiColors[outAnsiColorCounter%len(outAnsiColors)]
	outAnsiColorCounter++
	return ansiColor
}

func getErrAnsiColor() anzi.Color {
	errAnsiColorMutex.Lock()
	defer errAnsiColorMutex.Unlock()
	ansiColor := errAnsiColors[errAnsiColorCounter%len(errAnsiColors)]
	errAnsiColorCounter++
	return ansiColor
}
