package target

type Describer interface {
	Describe(Service) ([]string, error)
}

type Upper interface {
	Up(Service) error
}

type Downer interface {
	Down(Service) error
}

type Puller interface {
	Pull(Service) error
}

type Runner interface {
	Run(Service) error
}

type Driver {
	Describer
	Upper
	Downer
	//Puller
	//Runer
}
