package migration

type MockDriver struct {
	migrationTableExists func() (bool, error)
	createMigrationTable func() error
	listAppliedMigrations func() ([]uint64, error)
	startTransaction func() (Transaction, error)
}

func (d MockDriver) MigrationTableExists() (bool, error) {
	return d.migrationTableExists()
}

func (d MockDriver) CreateMigrationTable() error {
	return d.createMigrationTable()
}

func (d MockDriver) ListAppliedMigrations() ([]uint64, error) {
	return d.listAppliedMigrations()
}

func (d MockDriver) StartTransaction() (Transaction, error) {
	return d.startTransaction()
}

type MockReader struct {
	seek func(int64, int) (int64, error)
	read func(p []byte) (int, error)
}

func (r *MockReader) Seek(offset int64, whence int) (int64, error) {
	return r.seek(offset, whence)
}

func (r *MockReader) Read(p []byte) (int, error) {
	return r.read(p)
}
