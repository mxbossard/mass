package target

import "mby.fr/mass/internal/resources"

type Describer interface {
	Describe(resources.Service) ([]string, error)
}

type Upper interface {
	Up(resources.Service) error
}

type Downer interface {
	Down(resources.Service) error
}

type Puller interface {
	Pull(resources.Service) error
}

type Runner interface {
	Run(resources.Service) error
}

type Driver struct {
	Describer
	Upper
	Downer
	//Puller
	//Runer
}
