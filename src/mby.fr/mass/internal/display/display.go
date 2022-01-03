package display

import (
	"os"
	//"fmt"
	//"strings"
	"log"

	//"mby.fr/mass/internal/config"
	//"mby.fr/mass/internal/templates"
)

type Stringer interface {
	String() string
}

type Tuner func(string) string

var (
	logTuner Tuner = func(s string) string {
		return s
	}
	debugTuner = logTuner
	traceTuner = logTuner
	warnTuner = logTuner
)

func tune(tuners *[]Tuner, s string) string {
	if s == "" {
		return ""
	}
	for _, tuner := range *tuners {
		res := tuner(s)
		if res != "" {
			return res
		}
	}
	return ""
}

type ActionLogger struct {
	level uint8
	kind string
	tuners *[]Tuner
	printer *Printer
}

func (a ActionLogger) Nested(kind string) ActionLogger {
	return ActionLogger{a.level + 1, kind, a.tuners, a.printer}
}

func (a ActionLogger) Start() {

}

func (a ActionLogger) Progress(percent int8) {

}

func (a ActionLogger) End() {

}

func printMessages(printerFunc func(...interface{}) error, tuners *[]Tuner, defaultTuner Tuner, stringers ...Stringer) {
	strings := make([]interface{}, 5)
	for _, s := range stringers {
		allTuners := append(*tuners, defaultTuner)
		tuned := tune(&allTuners, s.String())
		strings = append(strings, tuned)
	}
	err := printerFunc(strings...)
	if err != nil {
		log.Fatal(err)
	}
}

func (a ActionLogger) Log(stringers ...Stringer) {
	printMessages((*a.printer).Out, a.tuners, logTuner, stringers...)
}

func (a ActionLogger) Debug(stringers ...Stringer) {
	printMessages((*a.printer).Out, a.tuners, debugTuner, stringers...)

}

func (a ActionLogger) Trace(stringers ...Stringer) {
	printMessages((*a.printer).Out, a.tuners, traceTuner, stringers...)
}

func (a ActionLogger) Warn(stringers ...Stringer) {
	printMessages((*a.printer).Out, a.tuners, warnTuner, stringers...)
}

func (a ActionLogger) Error(stringers ...Stringer) {

}

func (a ActionLogger) Fatal(stringers ...Stringer) {

}

type Displayer interface {
	Display(...interface{})
	ActionLogger(string) ActionLogger
}

type StandarDisplay struct {
	printer *Printer
	tuners *[]Tuner
}

func (d StandarDisplay) Display(objects ...interface{}) {
	err := (*d.printer).Print(objects...)
	if err != nil {
		log.Fatal(err)
	}
}

func (d StandarDisplay) ActionLogger(kind string) ActionLogger {
	return ActionLogger{0, kind, d.tuners, d.printer}
}

func New() Displayer {
	var printer Printer = Basic{os.Stdout, os.Stderr}
	tuners := []Tuner{}
	return StandarDisplay{&printer, &tuners}
}

