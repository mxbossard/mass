package display

import (
	"io"
	"fmt"
	"reflect"
	//"strings"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/templates"
	"mby.fr/utils/errorz"
)

// TODO: Pass a writer to the printer and unittest it with a string writer.

type Printer interface {
	Out(...interface{}) error
	Err(...interface{}) error
	Print(...interface{}) error
}

type Basic struct {
	out, err io.Writer
}

func (p Basic) Out(objects ...interface{}) (err error) {
	_, err = fmt.Fprint(p.out, objects...)
	return
}

func (p Basic) Err(objects ...interface{}) (err error) {
	_, err = fmt.Fprint(p.err, objects...)
	return
}

func (p Basic) Print(objects ...interface{}) (err error) {
	for _, obj := range objects {

		// Recursive call if obj is an array or a slice
		t := reflect.TypeOf(obj)
		if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
			arrayValue := reflect.ValueOf(obj)
			for i := 0; i < arrayValue.Len(); i++ {
				value := arrayValue.Index(i).Interface()
				err = p.Print(value)
				if err != nil {
					return
				}
			}
			continue
		}

		switch o:= obj.(type) {
		case errorz.Aggregated:
			for _, err := range o.Errors() {
				fmt.Fprintf(p.err, "Error: %s !\n", err)
			}
		case error:
			fmt.Fprintf(p.err, "Error: %s !\n", err)
		case config.Config:
			renderer := templates.New("")
			renderer.Render("display/basic/config.tpl", p.out, o)
		default:
			err = fmt.Errorf("Unable to Print object of type: %T", obj)
			return
		}
	}
	return
}
