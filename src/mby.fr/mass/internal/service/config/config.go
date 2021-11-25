package config

import (
	"fmt"
	"sync"
)

type MassConfig struct {
	workspacePath string
}

func (c MassConfig) String() string {
	return fmt.Sprintf("MassConfig workspacePath: %s.", c.workspacePath)
}



var lock = &sync.Mutex{}

type MassConfigService struct {
	config *MassConfig
}

func (s MassConfigService) Config() *MassConfig {
	return s.config
}

func newMassConfigService(workspacePath string) MassConfigService {
	maasConfig := &MassConfig{workspacePath}
	configService := MassConfigService{maasConfig}
	return configService
}

var maasConfigService *MassConfigService

func GetMassConfigService() *MassConfigService {
	if maasConfigService == nil {
		lock.Lock()
		defer lock.Unlock()
		if maasConfigService == nil {
			service := newMassConfigService()
			maasConfigService = &service
		}
	}
	return maasConfigService
}


