package display

import (
	//"os"
	//"fmt"
	//"strings"
	"log"

	//"mby.fr/mass/internal/config"
	//"mby.fr/mass/internal/templates"
)

type Displayer interface {
	Display(...interface{})
}

type StandarDisplay struct {
	printer Printer
}

func (d StandarDisplay) Display(objects ...interface{}) {
	err := d.printer.Print(objects...)
	if err != nil {
		log.Fatal(err)
	}
}

func New() Displayer {
	printer := Basic{}
	return StandarDisplay{printer}
}

