package workspace

import (
	"log"
	"os"

	"mby.fr/mass/internal/service/config"
)

func InitWorkspace(name string) {
	err := os.Mkdir(name, 0755)
	if (err != nil) {
		log.Fatal(err)
	}

	err = os.Chdir(name)
	if (err != nil) {
		log.Fatal(err)
	}

	config.InitMassConfig()
}

