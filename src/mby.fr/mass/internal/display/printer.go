package display

import (
	"os"
	"fmt"
	//"strings"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/templates"
	"mby.fr/utils/errorz"
)

// TODO: Pass a writer to the printer and unittest it with a string writer.

type Printer interface {
	//Error(error)
	//Config(config.Config)
	//AggregatedError(errorz.Aggregated)
	Print(...interface{}) error
}

type Basic struct {

}

func (p Basic) Error(err error) {
	fmt.Printf("Error: %s !\n", err)
}

func (p Basic) Config(c config.Config) {
	renderer := templates.New("")
	renderer.Render("display/basic/config.tpl", os.Stdout, c)
}

func (p Basic) AggregatedError(errors errorz.Aggregated) {
	for _, err := range errors.Errors() {
		p.Error(err)
	}
}


func (p Basic) Print(objects ...interface{}) (err error) {
	for _, obj := range objects {
		switch o:= obj.(type) {
		case errorz.Aggregated:
			p.AggregatedError(o)
		case error:
			p.Error(o)
		case config.Config:
			p.Config(o)
		case []config.Config:
			//p.Print(o...)
			for _, c:= range o {
				p.Print(c)
			}
		default:
			err = fmt.Errorf("Unable to Print object of type: %T", obj)
			return
		}
	}
	return
}
