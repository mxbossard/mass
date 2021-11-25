package config

import (
	"fmt"
	"sync"
)

type MaasConfig struct {
	workspacePath string
}

func (c MaasConfig) String() string {
	return fmt.Sprintf("MaasConfig workspacePath: %s.", c.workspacePath)
}



var lock = &sync.Mutex{}

type MaasConfigService struct {
	config *MaasConfig
}

func (s MaasConfigService) Config() *MaasConfig {
	return s.config
}

func makeMaasConfigService() *MaasConfigService {
	maasConfig := &MaasConfig{"badPath"}
	configService := &MaasConfigService{maasConfig}
	return configService
}

var maasConfigService *MaasConfigService

func GetMaasConfigService() *MaasConfigService {
	if maasConfigService == nil {
		lock.Lock()
		defer lock.Unlock()
		if maasConfigService == nil {
			maasConfigService = makeMaasConfigService()
		}
	}
	return maasConfigService
}


