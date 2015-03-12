package driver

// Conn is the interface used by the program to manipulate the migration table.
type Conn interface {
	HasTable() (bool, error)
	CreateTable() error
	GetApplied() ([]string, error)
	AddApplied(string) error
	RemoveApplied(string) error
	Execute(string) error
}
