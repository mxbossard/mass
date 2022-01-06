package display

import (
	//"os"
	"fmt"
	//"strings"
	"log"

	//"mby.fr/mass/internal/config"
	//"mby.fr/mass/internal/templates"
)

const ansiClear = "\033[0m"
const ansiBlack = "\033[1;30m"

const traceAnsiColor = "\033[0;97;46m"
const debugAnsiColor = "\033[0;97;44m"
const infoAnsiColor = "\033[0;97;42m"
const warnAnsiColor = "\033[0;31;43m"
const errorAnsiColor = "\033[0;97;41m"
const fatalAnsiColor = "\033[0;97;45m"

var okAnsiColors []string = []string{ansiClear, "\033[0;91m", "\033[0;92m", "\033[0;93m", "\033[0;94m", "\033[0;95m", "\033[0;96m", "\033[0;31m", "\033[0;32m", "\033[0;33m", "\033[0;34m", "\033[0;35m", "\033[0;36m"}

type Stringer interface {
	String() string
}

type Tuner func(interface{}) interface{}

var (
	logTuner Tuner = func(msg interface{}) interface{} {
		return msg
	}
	warnTuner = logTuner
	infoTuner = logTuner
	debugTuner = logTuner
	traceTuner = logTuner
)

func formatTuner(format string) Tuner {
	return func(msg interface{}) interface{} {
                return Format(format, msg)
        }
}

func tune(tuners *[]Tuner, msg interface{}) interface{} {
	if msg == nil {
		return nil
	}
	for _, tuner := range *tuners {
		msg = tuner(msg)
	}
	return msg
}

type ansiFormatted struct {
	format string
	content interface{}
	tab bool
	leftPad, rightPad int
}

//func (f ansiFormatted) String() string {
//	return fmt.Sprintf("%s%s%s", f.format, f.content, ansiClear)
//}

func Format(format string, object interface{}) ansiFormatted {
	return ansiFormatted{format: format, content: object}
}

var colorCounter = 0
func getOkAnsiColor() string {
	ansiColor := okAnsiColors[colorCounter % len(okAnsiColors)]
	colorCounter ++
	return ansiColor
}

type ActionLogger struct {
	level uint8
	action, subject string
	tuners *[]Tuner
	printer *Printer
	ansiColor string
}

func (a ActionLogger) Nested(action, subject string) ActionLogger {
	var printer Printer = NewPrinter(NewStandardOutputs())
	return ActionLogger{a.level + 1, action, subject, a.tuners, &printer, getOkAnsiColor()}
}

func (a ActionLogger) Start() {

}

func (a ActionLogger) Progress(percent int8) {

}

func (a ActionLogger) End() {

}

func formatLevel(format, level string) ansiFormatted {
	f := Format(format, fmt.Sprintf("[%s]", level))
	//f.leftPad = 7
	f.rightPad = 8
	return f
}

func printActionMessages(printerFunc func(...interface{}) error, actionLogger ActionLogger, level string, messages ...interface{}) {
	colorTuner := formatTuner(actionLogger.ansiColor)
	tuners := append(*actionLogger.tuners, colorTuner)

	// Log prefix
	prefix := Format("", fmt.Sprintf("%s(%s)", actionLogger.action, actionLogger.subject))
	prefix.rightPad = 25

	// Log level
	var formattedLevel interface{}
	switch level {
	case "":
		//formattedLevel = formatLevel("", "")
		formattedLevel = ""
	case "TRACE":
		formattedLevel = formatLevel(traceAnsiColor, level)
	case "DEBUG":
		formattedLevel = formatLevel(debugAnsiColor, level)
	case "INFO":
		formattedLevel = formatLevel(infoAnsiColor, level)
	case "WARN":
		formattedLevel = formatLevel(warnAnsiColor, level)
	case "ERROR":
		formattedLevel = formatLevel(errorAnsiColor, level)
	case "FATAL":
		formattedLevel = formatLevel(fatalAnsiColor, level)
	default:
		log.Fatalf("Unknown log level: %s", level)
	}

	allMessages := make([]interface{}, 0)
	allMessages = append(allMessages, prefix, formattedLevel)
	allMessages = append(allMessages, messages...)
	tunedMessages := make([]interface{}, 0)

	for _, msg := range allMessages {
		tuned := tune(&tuners, msg)
		if tuned != nil {
			tunedMessages = append(tunedMessages, tuned)
		} else {
			tunedMessages = append(tunedMessages, msg)
		}
	}
	tunedMessages = append(tunedMessages, "\n")
	err := printerFunc(tunedMessages...)
	if err != nil {
		log.Fatal(err)
	}
}

func (a ActionLogger) Log(messages ...interface{}) {
	printActionMessages((*a.printer).Out, a, "", messages...)
}

func (a ActionLogger) Info(messages ...interface{}) {
	printActionMessages((*a.printer).Out, a, "INFO", messages...)

}

func (a ActionLogger) Debug(messages ...interface{}) {
	printActionMessages((*a.printer).Out, a, "DEBUG", messages...)

}

func (a ActionLogger) Trace(messages ...interface{}) {
	printActionMessages((*a.printer).Out, a, "TRACE", messages...)
}

func (a ActionLogger) Warn(messages ...interface{}) {
	printActionMessages((*a.printer).Out, a, "WARN", messages...)
}

func (a ActionLogger) Error(messages ...interface{}) {
	printActionMessages((*a.printer).Err, a, "ERROR", messages...)
}

func (a ActionLogger) Fatal(messages ...interface{}) {
	printActionMessages((*a.printer).Err, a, "FATAL", messages...)

}

type Displayer interface {
	Display(...interface{})
	ActionLogger(string, string) ActionLogger
}

type StandarDisplay struct {
	printer *Printer
	tuners *[]Tuner
}

func (d StandarDisplay) Display(objects ...interface{}) {
	err := (*d.printer).Out(objects...)
	if err != nil {
		log.Fatal(err)
	}
}

func (d StandarDisplay) ActionLogger(action, subject string) ActionLogger {
	var printer Printer = NewPrinter(NewStandardOutputs())
	return ActionLogger{0, action, subject, d.tuners, &printer, getOkAnsiColor()}
}

func New() Displayer {
	var printer Printer = NewPrinter(NewStandardOutputs())
	tuners := []Tuner{}
	return StandarDisplay{&printer, &tuners}
}

