package display

import (
	"os"
	//"fmt"
	//"strings"

	"mby.fr/mass/internal/config"
	"mby.fr/mass/internal/templates"
)

type Display interface {
	Config(config.Config)
}

type Basic struct {

}

func (d Basic) Config(c config.Config) {
	renderer := templates.New("")
	renderer.Render("display/basic/config.tpl", os.Stdout, c)
}

