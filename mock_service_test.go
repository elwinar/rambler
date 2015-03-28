package main

type MockService struct {
	initialized func() (bool, error)
	initialize  func() error
	available   func() ([]*Migration, error)
	applied     func() ([]*Migration, error)
	apply       func(*Migration) error
	reverse     func(*Migration) error
}

func (s MockService) Initialized() (bool, error) {
	return s.initialized()
}

func (s MockService) Initialize() error {
	return s.initialize()
}

func (s MockService) Available() ([]*Migration, error) {
	return s.available()
}

func (s MockService) Applied() ([]*Migration, error) {
	return s.applied()
}

func (s MockService) Apply(migration *Migration) error {
	return s.apply(migration)
}

func (s MockService) Reverse(migration *Migration) error {
	return s.reverse(migration)
}
