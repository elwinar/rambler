package main

import (
	"github.com/gonuts/commander"
	"github.com/gonuts/flag"
	"os"
)

var (
	// The database flagset. There shouldn't be another flagset unless the command
	// grow heavily in functionnalities.
	flagsetDatabase = flag.NewFlagSet("database", flag.ExitOnError)

	// The global command that handle the program
	command = &commander.Command{
		UsageLine: os.Args[0] + " COMMAND [options]",
		Short:     "Manipulate databases schemas",
		Subcommands: []*commander.Command{
			commandCreate,
			commandMigrate,
			commandRollback,
		},
	}

	// The migrate command
	commandMigrate = &commander.Command{
		Run:       migrate,
		UsageLine: "migrate [options]",
		Short:     "Apply migrations on a database",
	}

	// The rollback command
	commandRollback = &commander.Command{
		Run:       rollback,
		UsageLine: "rollback [options]",
		Short:     "Rollback migrations on a database",
	}
	
	// The create command
	commandCreate = &commander.Command{
		Run:       create,
		UsageLine: "create DESCRIPTION",
		Short:     "Create a new migration file",
	}
)

// Load the various flags needed by subcommands
func init() {
	// The flagset database is probably the only one flagset needed
	flagsetDatabase.String("protocol", "tcp", "the protocol to communicate with the database")
	flagsetDatabase.String("host", "localhost", "the server which host the database")
	flagsetDatabase.Int("port", 3306, "the port on which the database listen")
	flagsetDatabase.String("user", "root", "the user to connect to the database server")
	flagsetDatabase.String("password", "", "the password of the user on the database server")
	flagsetDatabase.String("db", "", "the name of the database (required)")

	// Migrate and rollback both use the database flagset
	commandMigrate.Flag = *flagsetDatabase
	commandRollback.Flag = *flagsetDatabase
}
