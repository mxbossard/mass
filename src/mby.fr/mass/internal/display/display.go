package display

import (
	"fmt"
	"log"
	"sync"
	"time"

	"mby.fr/mass/internal/logger"
	"mby.fr/mass/internal/output"
	"mby.fr/utils/logz"
)

// Display should display data to user.
// Displayed data types :
// Logs
// Stdout & Stderr
// Errors
// Printed objects

const (
	logPeriodInSeconds = 5
)

type Displayer interface {
	logz.Logger
	Display(...interface{})
	BufferedActionLogger(string, string) logger.ActionLogger
	ImmediateActionLogger(string, string) logger.ActionLogger
	Flush() error
}

type StandarDisplay struct {
	sync.Mutex
	*logger.ActionLogger
	outs              output.Outputs
	printer           *Printer
	tuners            *[]Tuner
	flushableOuts     *[]output.Outputs // FIXME memory leak this slice is never cleaned
	bufferedOuts      *[]output.Outputs // FIXME memory leak this slice is never cleaned
	mainOuts          output.Outputs
	lastMainOutsWrite time.Time
}

func (d StandarDisplay) Display(objects ...interface{}) {
	err := (*d.printer).Out(objects...)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *StandarDisplay) BufferedActionLogger(action, subject string) logger.ActionLogger {
	return actionLogger(d, action, subject, true)
}

func (d *StandarDisplay) ImmediateActionLogger(action, subject string) logger.ActionLogger {
	return actionLogger(d, action, subject, false)
}

func (d StandarDisplay) Flush() (err error) {
	err = (*d.printer).Flush()
	if err != nil {
		return err
	}

	for _, outs := range *d.flushableOuts {
		fmt.Printf("Flushing out ...\n")
		err = outs.Flush()
		if err != nil {
			return err
		}
	}

	d.Debug("Flushing display.")
	return
}

func (d *StandarDisplay) flushMainOutputs() {
	for {
		d.Lock()
		if time.Now().Sub(d.lastMainOutsWrite).Seconds() > logPeriodInSeconds {
			// If main outputs did not write for 5 seconds select new main outputs
			for _, outs := range *d.bufferedOuts {
				//fmt.Println("outs last write:", d.mainOuts, outs.LastWriteTime(), time.Now().Sub(d.lastMainOutsWrite).Seconds(), time.Now().Sub(outs.LastWriteTime()).Seconds())
				if time.Now().Sub(outs.LastWriteTime()).Seconds() < logPeriodInSeconds {
					d.mainOuts = outs
					//break
				}
			}
		}
		if d.mainOuts != nil {
			//fmt.Println("Flushing main printer")
			d.mainOuts.Flush()
			d.lastMainOutsWrite = d.mainOuts.LastWriteTime()
		}
		d.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func (d *StandarDisplay) flushOtherOutputs() {
	for {
		time.Sleep(logPeriodInSeconds * time.Second)
		d.Lock()
		if len(*d.bufferedOuts) > 0 {
			for _, outs := range *d.bufferedOuts {
				if outs != d.mainOuts {
					outs.Flush()
				}
			}
		}
		d.Unlock()
	}
}

func actionLogger(d *StandarDisplay, action, subject string, buffered bool) logger.ActionLogger {
	var outs output.Outputs
	if buffered {
		outs = output.NewBufferedOutputs(d.outs)
		appended := append(*d.bufferedOuts, outs)
		d.bufferedOuts = &appended
	} else {
		outs = d.outs
	}
	al := logger.NewAction(outs, action, subject)
	appended := append(*d.flushableOuts, outs)
	d.flushableOuts = &appended
	return al
}

func newInstance() Displayer {
	var m sync.Mutex
	outs := output.NewStandardOutputs()
	var printer Printer = NewPrinter(NewStandardOutputs())
	tuners := []Tuner{}
	flushableOuts := []output.Outputs{}
	bufferedOuts := []output.Outputs{}
	d := StandarDisplay{m, nil, outs, &printer, &tuners, &flushableOuts, &bufferedOuts, nil, time.Time{}}
	logger := actionLogger(&d, "display", "", false)
	d.ActionLogger = &logger

	go d.flushMainOutputs()
	go d.flushOtherOutputs()

	return &d
}

var service = newInstance()

func Service() Displayer {
	return service
}

// Tune objects converting them to other objects
type Tuner func(interface{}) interface{}
