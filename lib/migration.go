package lib

import (
	"bufio"
	"errors"
	"github.com/elwinar/viper"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	prefix = `-- rambler`
)

// Migration represent a schema migration file, composed of up and down queries
type Migration struct {
	Version     uint64    `db:"version"`
	Date        time.Time `db:"date"`
	Description string    `db:"description"`
	File        string    `db:"file"`
}

// NewMigration create a migration struct with the given path
func NewMigration(path string) (Migration, error) {
	m := Migration{}

	if !IsMigrationFile(path) {
		return m, errors.New("the specified path isn't a migration file")
	}

	chunks := strings.SplitN(path, "_", 2)
	m.Version, _ = strconv.ParseUint(chunks[0], 10, 64)
	m.Date = time.Now()
	m.Description = strings.Replace(strings.TrimSuffix(chunks[1], ".sql"), "_", " ", -1)
	m.File = path

	return m, nil
}

// IsMigrationFile check whether the given path follow the migration file naming convention
func IsMigrationFile(path string) bool {
	match, err := regexp.MatchString(`([0-9]+)_([a-zA-Z0-9_-]+).sql`, path)
	if err != nil {
		panic(err)
	}
	if !match {
		return false
	}
	return true
}

// GetMigrationsDir returns the path to the migration directory
func GetMigrationsDir() string {
	return filepath.Join(filepath.Dir(viper.ConfigFileUsed()), viper.GetString("migrations"))
}

// Scan open the migration file and parse it line by line, keeping only lines in the
// section passed as parameter.
func (m Migration) Scan(section string) ([]string, error) {
	file, err := os.Open(path.Join(viper.GetString("migrations"), m.File))
	if err != nil {
		return nil, err
	}
	var scanner = bufio.NewScanner(file)
	var statements []string
	var buffer string
	recording := false
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, prefix) {
			if len(strings.TrimSpace(buffer)) != 0 {
				statements = append(statements, strings.TrimSpace(buffer))
			}
			buffer = ""
			cmd := strings.TrimSpace(line[len(prefix):])
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
		statements = append(statements, strings.TrimSpace(buffer))
	}
	return statements, nil
}

func (m Migration) IsAvailable() bool {
	_, err := os.Stat(m.File)
	return os.IsNotExist(err)
}
