package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bradfitz/slice"
	"github.com/elwinar/rambler/driver"
	_ "github.com/elwinar/rambler/driver/mysql"
	_ "github.com/elwinar/rambler/driver/postgresql"
	_ "github.com/elwinar/rambler/driver/sqlite"
)

var (
	// ErrNilMigration is returned when the service is given a nil migration
	ErrNilMigration = errors.New("nil migration")
)

// Service is the struct that gather operations to manipulate the
// database and migrations on disk
type Service struct {
	conn   driver.Conn
	env    Environment
	dryRun bool
}

// NewService initialize a new service with the given environment
func NewService(env Environment, dryRun bool) (*Service, error) {
	fi, err := os.Stat(env.Directory)
	if err != nil {
		return nil, fmt.Errorf("directory %s unavailable: %s", env.Directory, err.Error())
	}

	if !fi.Mode().IsDir() {
		return nil, fmt.Errorf("%s isn't a directory", env.Directory)
	}

	driver, err := driver.Get(env.Driver)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize driver: %s", err.Error())
	}

	conn, err := driver.New(env.Config)
	if err != nil {
		return nil, fmt.Errorf("unable to initialize connection: %s", err.Error())
	}

	return &Service{
		conn:   conn,
		env:    env,
		dryRun: dryRun,
	}, nil
}

// Initialized check if the migration table exists in the
// database
func (s Service) Initialized() (bool, error) {
	return s.conn.HasTable()
}

// Initialize create the migration table in the database
func (s Service) Initialize() error {
	return s.conn.CreateTable()
}

// Available return the migrations in the environment's directory sorted in
// ascending lexicographic order.
func (s Service) Available() ([]*Migration, error) {
	files, _ := filepath.Glob(filepath.Join(s.env.Directory, "*.sql")) // The only possible error here is a pattern error

	var migrations []*Migration
	for _, file := range files {
		migration, err := NewMigration(file)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	slice.Sort(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

// Applied return the migrations in the environment's directory that are marked
// as applied in the database sorted in ascending lexicographic order.
func (s Service) Applied() ([]*Migration, error) {
	files, err := s.conn.GetApplied()
	if err != nil {
		return nil, err
	}

	var migrations []*Migration
	for _, file := range files {
		migration, err := NewMigration(filepath.Join(s.env.Directory, file))
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migration)
	}

	slice.Sort(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	return migrations, nil
}

// Apply execute the up statements of the given migration to the
// database then mark the migration as applied
func (s Service) Apply(migration *Migration, save bool) error {
	if migration == nil {
		return ErrNilMigration
	}

	for _, statement := range migration.Up() {
		if s.dryRun {
			logger.Info("statement: %s", statement)
			continue
		}

		err := s.conn.Execute(statement)
		if err != nil {
			return fmt.Errorf("unable to apply migration %s: %s\n%s", migration.Name, err, statement)
		}
	}

	if s.dryRun {
		return nil
	}

	if !save {
		return nil
	}

	err := s.conn.AddApplied(migration.Name)
	if err != nil {
		return fmt.Errorf("unable to mark migration %s as applied: %s", migration.Name, err)
	}

	return nil
}

// Reverse execute the down statements of the given migration to the
// database then mark the migration as not applied
func (s Service) Reverse(migration *Migration, save bool) error {
	if migration == nil {
		return ErrNilMigration
	}

	for _, statement := range migration.Down() {
		if s.dryRun {
			logger.Info("statement: %s", statement)
			continue
		}

		err := s.conn.Execute(statement)
		if err != nil {
			return fmt.Errorf("unable to reverse migration %s: %s\n%s", migration.Name, err, statement)
		}
	}

	if s.dryRun {
		return nil
	}

	if !save {
		return nil
	}

	err := s.conn.RemoveApplied(migration.Name)
	if err != nil {
		return fmt.Errorf("unable to mark migration %s as not applied: %s", migration.Name, err)
	}

	return nil
}
