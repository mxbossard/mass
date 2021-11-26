package workspace

import (
	"log"
	"os"
)

func Init(name string) {
	err := os.Mkdir(name, 0755)
	if (err != nil) {
		log.Fatal(err)
	}

	err = os.Chdir(name)
	if (err != nil) {
		log.Fatal(err)
	}

	InitConfig()
}

