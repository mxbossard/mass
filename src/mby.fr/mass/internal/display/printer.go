package display

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/templates"
	"mby.fr/utils/ansi"
	"mby.fr/utils/errorz"
)

// Outputs responsible for keeping reference of outputs writers (example: stdout, file, ...)
// Printer responsible for printing messages in outputs (example: print with colors, without colors, ...)

type Flusher interface {
	Flush() error
}

type Outputs interface {
	Flusher
	Out() io.Writer
	Err() io.Writer
}

type BasicOutputs struct {
	out, err io.Writer
}

func (o BasicOutputs) Flush() error {
	outs := []io.Writer{o.out, o.err}
	for _, out := range outs {
		f, ok := out.(Flusher)
		if ok {
			err := f.Flush()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (o BasicOutputs) Out() io.Writer {
	return o.out
}

func (o BasicOutputs) Err() io.Writer {
	return o.err
}

type Printer interface {
	Outputs() Outputs
	Flush() error
	Out(...interface{}) error
	Err(...interface{}) error
	//Print(...interface{}) error
	LastPrint() time.Time
}

type BasicPrinter struct {
	sync.Mutex
	outputs   Outputs
	lastPrint time.Time
}

func (p *BasicPrinter) Outputs() Outputs {
	return p.outputs
}

func (o BasicPrinter) Flush() error {
	o.Lock()
	defer o.Unlock()
	return o.outputs.Flush()
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
	switch o := obj.(type) {
	case string:
		str = o
	case int:
		str = strconv.Itoa(o)
	case float64:
		str = strconv.FormatFloat(o, 'E', 3, 32)

	case ansiFormatted:
		if o.content == "" {
			return "", nil
		}
		content, err := stringify(o.content)
		if err != nil {
			return "", err
		}
		if o.format != "" {
			str = fmt.Sprintf("%s%s%s", o.format, content, ansi.Reset)
		} else {
			str = content
		}
		if o.tab {
			str += "\t"
		} else if o.leftPad > 0 {
			spaceCount := o.leftPad - len(content)
			if spaceCount > 0 {
				str = strings.Repeat(" ", spaceCount) + str
			}
		} else if o.rightPad > 0 {
			spaceCount := o.rightPad - len(content)
			if spaceCount > 0 {
				str += strings.Repeat(" ", spaceCount)
			}
		}
	case error:
		str = fmt.Sprintf("Error: %s !\n", obj)
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

func expandObjects(objects ...interface{}) (allObjects []interface{}) {
	for _, obj := range objects {
		// Recursive call if obj is an array or a slice
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			arrayValue := reflect.ValueOf(obj)
			for i := 0; i < arrayValue.Len(); i++ {
				value := arrayValue.Index(i).Interface()
				expanded := expandObjects(value)
				allObjects = append(allObjects, expanded...)
			}
			continue
		} else {
			allObjects = append(allObjects, obj)
		}
	}
	return
}

func printTo(w io.Writer, objects ...interface{}) (err error) {
	objects = expandObjects(objects)
	var toPrint []string

	for _, obj := range objects {
		switch o := obj.(type) {
		case errorz.Aggregated:
			for _, err := range o.Errors() {
				//fmt.Fprintf(w, "Error: %s !\n", err)
				printTo(w, err)
			}
		default:
			str, err := stringify(o)
			if err != nil {
				return err
			}
			if len(str) > 0 {
				//fmt.Printf("adding str: %d [%s](%T)\n", len(str), str, str)
				toPrint = append(toPrint, str)
			}
		}
	}
	//fmt.Printf("toPrint: %d %s\n", len(toPrint), toPrint)
	_, err = fmt.Fprintf(w, "%s", strings.Join(toPrint, " "))
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

// ANSI formatting for content
type ansiFormatted struct {
	format            string
	content           interface{}
	tab               bool
	leftPad, rightPad int
}
