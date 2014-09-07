package main

import (
	"fmt"
	"github.com/gonuts/commander"
	"io/ioutil"
	"strings"
	"time"
)

const (
	fmtDate      = `20060102150405`
	tplFilename  = `%s_%s.sql`
	tplMigration = `-- +rambler Up

-- +rambler Down

`
)

func create(command *commander.Command, args []string) error {
	ioutil.WriteFile(fmt.Sprintf(tplFilename, time.Now().Format(fmtDate), strings.Replace(args[0], " ", "_", -1)), []byte(tplMigration), 0644)
	return nil
}
