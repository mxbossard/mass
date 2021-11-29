package workspace

import (
)

const configDir = "config"

func InitConfig() {
	workspacePath := GetWorkDirPath()
	CreateNewSubDirectory(workspacePath, configDir)
}

