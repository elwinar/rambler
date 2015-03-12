package main

type MockConn struct {
	hasTable func() (bool, error)
	createTable func() error
	getApplied func() ([]string, error)
	addApplied func(string) error
	removeApplied func(string) error
	execute func(string) error
}

func (c MockConn) HasTable() (bool, error) {
	return c.hasTable()
}

func (c MockConn) CreateTable() error {
	return c.createTable()
}

func (c MockConn) GetApplied() ([]string, error) {
	return c.getApplied()
}

func (c MockConn) AddApplied(migration string) error {
	return c.addApplied(migration)
}

func (c MockConn) RemoveApplied(migration string) error {
	return c.removeApplied(migration)
}

func (c MockConn) Execute(statement string) error {
	return c.execute(statement)
}
