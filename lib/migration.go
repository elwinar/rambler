package lib

import (
	"bufio"
	"errors"
	"github.com/spf13/viper"
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

type Migration struct {
	Version     uint64    `db:"version"`
	Date        time.Time `db:"date"`
	Description string    `db:"description"`
	File        string    `db:"file"`
}

func NewMigration (path string) (Migration, error) {
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

func IsMigrationFile (path string) bool {
	match, err := regexp.MatchString(`([0-9]+)_([a-zA-Z0-9_-]+).sql`, path)
	if err != nil {
		panic(err)
	}
	if !match {
		return false
	}
	return true
}

func GetMigrationsDir () string {
	return filepath.Join(filepath.Dir(viper.ConfigFileUsed()), viper.GetString("migrations"))
}

func GetMigrationsFiles () ([]Migration, error) {
	files, err := filepath.Glob(GetMigrationsDir() + "/*.sql")
	if err != nil {
		return nil, err
	}
	
	var migrations []Migration
	for _, file := range files {
		migration, err := NewMigration(filepath.Base(file))
		if err != nil {
			continue
		}
		migrations = append(migrations, migration)
	}
	
	return migrations, nil
}

func GetMigrationsRows () ([]Migration, error) {
	var rows []Migration = nil
	err := db.Select(&rows, `SELECT * FROM migrations`)
	return rows, err
}

// scan open the migration file and parse it line by line, keeping only lines in the
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

func (m Migration) Up () ([]string, error) {
	return m.Scan("up")
}

func (m Migration) Down () ([]string, error) {
	sections, err := m.Scan("down")
	if err != nil {
		return sections, err
	}
	var reversed []string = make([]string, len(sections))
	for i := 0; i < len(sections); i++ {
		reversed[len(sections) - i - 1] = sections[i]
	}
	return reversed, nil
}
