package workspace

import (
	"log"
	"os"
)

var defaultEnvs = []string{"dev", "stage", "prod"}

func Init(name string) {
	CreateNewDirectory(".", name)

        err := os.Chdir(name)
        if (err != nil) {
                log.Fatal(err)
        }

	InitSettings()
	InitConfig()

	for _, envName := range defaultEnvs {
		InitEnv(envName)
	}
}

