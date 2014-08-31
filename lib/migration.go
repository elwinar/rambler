package lib

import (
	"errors"
	"github.com/spf13/viper"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
