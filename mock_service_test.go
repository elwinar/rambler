package main

type MockService struct {
	initialized func() (bool, error)
	initialize  func() error
	available   func() ([]Migration, error)
	applied     func() ([]Migration, error)
	apply       func(Migration, bool) error
	reverse     func(Migration, bool) error
}

func (s MockService) Initialized() (bool, error) {
	return s.initialized()
}

func (s MockService) Initialize() error {
	return s.initialize()
}

func (s MockService) Available() ([]Migration, error) {
	return s.available()
}

func (s MockService) Applied() ([]Migration, error) {
	return s.applied()
}

func (s MockService) Apply(migration Migration, save bool) error {
	return s.apply(migration, save)
}

func (s MockService) Reverse(migration Migration, save bool) error {
	return s.reverse(migration, save)
}
