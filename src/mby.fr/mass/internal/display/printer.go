package display

import (
	"io"
	"os"
	"bufio"
	"fmt"
	"reflect"
	"time"
	"sync"
	"log"
	"strings"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/templates"
	"mby.fr/utils/errorz"
)

// Outputs responsible for keeping reference of outputs writers (example: stdout, file, ...)
// Printer responsible for printing messages in outputs (example: print with colors, without colors, ...)

type Outputs interface {
	Out() io.Writer
	Err() io.Writer
}

type BasicOutputs struct {
	out, err io.Writer
}

func (o BasicOutputs) Out() io.Writer {
	return o.out
}

func (o BasicOutputs) Err() io.Writer {
	return o.err
}

type Printer interface {
	Out(...interface{}) error
	Err(...interface{}) error
	//Print(...interface{}) error
	LastPrint() time.Time
}

type BasicPrinter struct {
	sync.Mutex
	outputs Outputs
	lastPrint time.Time
}

func (p *BasicPrinter) Out(objects ...interface{}) (err error) {
	p.Lock()
	defer p.Unlock()
	p.lastPrint = time.Now()
	//_, err = fmt.Fprint(p.outputs.Out(), objects...)
	err = printTo(p.outputs.Out(), objects...)
	return
}

func (p *BasicPrinter) Err(objects ...interface{}) (err error) {
	p.Lock()
	defer p.Unlock()
	p.lastPrint = time.Now()
	//_, err = fmt.Fprint(p.outputs.Err(), objects...)
	err = printTo(p.outputs.Err(), objects...)
	return
}

func (p BasicPrinter) LastPrint() time.Time {
	return p.lastPrint
}

func stringify(obj interface{}) (str string, err error) {
	switch o:= obj.(type) {
	case string:
		str = o
	case ansiFormatted:
		content, err := stringify(o.content)
		if err != nil {
			return "", err
		}
		str = fmt.Sprintf("%s%s%s", o.format, content, ansiClear)
	case error:
		str = fmt.Sprintf("Error: %s !\n", err)
	case config.Config:
		renderer := templates.New("")
		builder := strings.Builder{}
		err = renderer.Render("display/basic/config.tpl", &builder, o)
		str = builder.String()
	default:
		err = fmt.Errorf("Unable to Print object of type: %T", obj)
		return
	}
	return
}

func printTo(w io.Writer, objects ...interface{}) (err error) {
	//p.Lock()
	//defer p.Unlock()
	//p.lastPrint = time.Now()
	var k = 0
	for _, obj := range objects {

		// Recursive call if obj is an array or a slice
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			arrayValue := reflect.ValueOf(obj)
			for i := 0; i < arrayValue.Len(); i++ {
				value := arrayValue.Index(i).Interface()
				err = printTo(w, value)
				if err != nil {
					return
				}
			}
			continue
		}

		var toPrint string

		switch o:= obj.(type) {
		//case string:
		//	toPrint = o
	        //case ansiFormatted:
		//	fmt.Fprintf(w, o.format)
	        //        printTo(w, o.content)
		//	fmt.Fprintf(w, ansiClear)
		case errorz.Aggregated:
			for _, err := range o.Errors() {
				//fmt.Fprintf(w, "Error: %s !\n", err)
				printTo(w, err)
			}
		//case error:
		//	toPrint = fmt.Sprintf("Error: %s !\n", err)
		//case config.Config:
		//	renderer := templates.New("")
		//	renderer.Render("display/basic/config.tpl", w, o)
		default:
		//	err = fmt.Errorf("Unable to Print object of type: %T", obj)
		//	return
			toPrint, err = stringify(o)
			if err != nil {
				return err
			}
		}

		if toPrint != "" {
			// k represent the printed objects count
			spacer := ""
			if k > 0 {
				spacer = ""
			}
			k ++
			//fmt.Printf("k: %d / obj(%T): %s / objects(%d): %T %s\n",  k, obj, obj, len(objects), objects)
			// Print space between printable objects
			_, err = fmt.Fprintf(w, "%s%s", toPrint, spacer)
			if err != nil {
				return
			}
		}
	}
	return
}

func NewStandardOutputs() Outputs {
	return BasicOutputs{os.Stdout, os.Stderr}
}

func NewBufferedOutputs(outputs Outputs) Outputs {
	buffOut := bufio.NewWriter(outputs.Out())
	buffErr := bufio.NewWriter(outputs.Err())
	buffered := BasicOutputs{buffOut, buffErr}
	return buffered
}

const logPeriodInSeconds = 5

var mainPrinter *BasicPrinter
var lastMainPrint time.Time
var flushablePrinters []*BasicPrinter

var globalMutex sync.Mutex

func flush(printer *BasicPrinter) {
	printer.Lock()
	defer printer.Unlock()
	err := printer.outputs.Out().(*bufio.Writer).Flush()
	if err != nil {
		log.Fatal(err)
	}
	err = printer.outputs.Err().(*bufio.Writer).Flush()
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	go flushMainPrinter()
	go flushOtherPrinters()
}

func flushMainPrinter() {
	for {
		globalMutex.Lock()
		if time.Now().Sub(lastMainPrint).Seconds() > logPeriodInSeconds {
			// If main printer did not print for 5 seconds select new main printer
			for _, printer := range flushablePrinters {
				//fmt.Println("printer last print:", printer.LastPrint(), time.Now().Sub(lastMainPrint).Seconds(), time.Now().Sub(printer.LastPrint()).Seconds())
				if time.Now().Sub(printer.LastPrint()).Seconds() < logPeriodInSeconds {
					mainPrinter = printer
					lastMainPrint = printer.LastPrint()
				}
			}
		}
		if mainPrinter != nil {
			//fmt.Println("Flushing main printer")
			flush(mainPrinter)
			lastMainPrint = mainPrinter.LastPrint()
		}
		globalMutex.Unlock()
		time.Sleep(100 * time.Millisecond)
	}
}

func flushOtherPrinters() {
	for {
		globalMutex.Lock()
		if len(flushablePrinters) > 0 {
			for _, printer := range flushablePrinters {
				if printer != mainPrinter {
					flush(printer)
				}
			}
		}
		globalMutex.Unlock()
		time.Sleep(logPeriodInSeconds * time.Second)
	}
}

func NewPrinter(outputs Outputs) Printer {
	buffered := NewBufferedOutputs(outputs)
	var m sync.Mutex
	var t time.Time
	printer := BasicPrinter{m, buffered, t}

	flushablePrinters = append(flushablePrinters, &printer)

	// Flush printer every 5 seconds
	//go func() {
	//	for {
	//		time.Sleep(5 * time.Second)
	//		flush(&printer)
	//	}
	//}()

	return &printer
}

