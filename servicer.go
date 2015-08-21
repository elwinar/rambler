package main

type Servicer interface {
	Initialized() (bool, error)
	Initialize() error
	Available() ([]*Migration, error)
	Applied() ([]*Migration, error)
	Apply(*Migration) error
	Reverse(*Migration) error
}
