package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	errInvalidName                   = `file %s is not a valid migration name`
	errMissingMigration              = `migration file %s is missing`
	errNoMigrationsTable             = `migrations table not found`
	errUnableToCreateMigrationsTable = `unable to create migrations table`
	queryCheckMigrationsTable        = `SELECT table_name
		FROM information_schema.tables 
		WHERE table_schema = '%s' 
		AND table_name = '%s'`
	queryCreateMigrationsTable = `CREATE TABLE migrations (
		version BIGINT UNSIGNED NOT NULL PRIMARY KEY,
		date DATETIME NOT NULL DEFAULT current_timestamp,
		description VARCHAR(255) NOT NULL,
		file VARCHAR(255) NOT NULL
	)`
	queryGetMigrations        = `SELECT * FROM %s ORDER BY version ASC`
	msgMigrationsTableCreated = `migrations table created`
	regexFilename             = `([0-9]+)_([a-zA-Z0-9_-]+).sql`
	sqlPrefix                 = `-- +rambler`
	tableMigrations           = `migrations`
)

// Migration represent a migration row from the migrations table.
// Each one should match a file in the migrations directory.
type Migration struct {
	Version     uint64    `db:"version"`
	Date        time.Time `db:"date"`
	Description string    `db:"description"`
	File        string    `db:"file"`
}

// Get the actual migrations from the database, and create the migrations table
// if it didn't exists.
func GetMigrations(db *sqlx.DB) ([]Migration, error) {
	// Check if the migrations table exists
	err := db.Get(new(struct {
		TableName string `db:"table_name"`
	}), fmt.Sprintf(queryCheckMigrationsTable, config.Database, tableMigrations))

	if err != nil {
		fmt.Println(errNoMigrationsTable)

		// Create the migrations table
		_, err := db.Exec(queryCreateMigrationsTable)
		if err != nil {
			return nil, errors.New(errUnableToCreateMigrationsTable + ": " + err.Error())
		}
		fmt.Println(msgMigrationsTableCreated)
	}

	// Get all the migrations. In case of error, stop here. Empty result set
	// shouldn't result in an error, as I used the Select method, which return
	// an empty slice if there is no result.
	var migrations []Migration
	err = db.Select(&migrations, fmt.Sprintf(queryGetMigrations, tableMigrations))
	if err != nil {
		return nil, err
	}

	return migrations, nil
}

// File holds all informations about a migration file
type MigrationFile struct {
	Version     uint64
	Description string
	File        string
}

// NewMigrationFile initialize a MigrationFile struct from a filename
func NewMigrationFile(path string) (*MigrationFile, error) {
	f := &MigrationFile{}

	regex := regexp.MustCompile(regexFilename)
	matches := regex.FindStringSubmatch(filepath.Base(path))
	if matches == nil {
		return nil, errors.New(fmt.Sprintf(errInvalidName, path))
	}
	version, _ := strconv.ParseUint(matches[1], 10, 64)

	f.Version = uint64(version)
	f.Description = strings.Replace(matches[2], "_", " ", -1)
	f.File = path

	return f, nil
}

// Up scan the migration file and returns the statements contained in Up sections
func (f MigrationFile) Up() ([]string, error) {
	return f.Scan("Up")
}

// Down scan the migration file and returns the statements contained in Down sections
func (f MigrationFile) Down() ([]string, error) {
	return f.Scan("Down")
}

// Scan open the migration file and parse it line by line, keeping only lines in the
// section passed as parameter.
func (f MigrationFile) Scan(section string) ([]string, error) {
	file, err := os.Open(f.File)
	if err != nil {
		return nil, err
	}

	var scanner = bufio.NewScanner(file)
	var statements []string
	var buffer string

	recording := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, sqlPrefix) {
			if len(strings.TrimSpace(buffer)) != 0 {
				statements = append(statements, buffer)
			}
			buffer = ""
			cmd := strings.TrimSpace(line[len(sqlPrefix):])
			switch cmd {
			case section:
				recording = true

			default:
				recording = false
			}
			continue
		}
		if recording {
			buffer = buffer + "\n" + line
		}
	}
	if len(strings.TrimSpace(buffer)) != 0 {
		statements = append(statements, buffer)
	}
	return statements, nil
}

// MigrationFiles is a wrapper for []MigrationFile that implement the sort.Interface
// interface.
type MigrationFiles []MigrationFile

// Len returns the number of migration files
func (slice MigrationFiles) Len() int {
	return len(slice)
}

// Less tells if the file i has a lower version number than file j
func (slice MigrationFiles) Less(i, j int) bool {
	return slice[i].Version < slice[j].Version
}

// Swap invert the position of the i and j files
func (slice MigrationFiles) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// GetMigrationFiles return all migration files that match the given pattern
func GetMigrationFiles(pattern string) (MigrationFiles, error) {
	var files MigrationFiles
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		file, err := NewMigrationFile(match)
		if err != nil {
			fmt.Println(err)
			continue
		}
		files = append(files, *file)
	}
	sort.Sort(files)
	return files, nil
}
