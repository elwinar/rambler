package commands

import (
	"github.com/elwinar/cobra"
	"github.com/elwinar/rambler/lib"
	"github.com/elwinar/viper"
	jww "github.com/spf13/jwalterweatherman"
	"sort"
)

func init() {
	// Add the command to the main command
	Rambler.AddCommand(Reverse)

	// Set the default configuration
	viper.SetDefault("all", false)

	// Add ubiquitous flags to the main command
	Reverse.Flags().BoolP("all", "a", false, "reverse all migrations")

	// Set overrides from the command-line to viper
	override("reverse-all", Reverse.Flags().Lookup("all"))
}

var Reverse = &cobra.Command{
	Use:   "reverse",
	Short: "Reverse the last migration",
	Run: do(func(cmd *cobra.Command, args []string) {
		// Start by opening the database connection and looking for the migration
		// table. If not found, we stop here.
		jww.TRACE.Println("Openning database connection")
		db, err := lib.GetDB()
		defer db.Close()
		if err != nil {
			jww.ERROR.Println("Unable to open database connection:", err)
			return
		}

		jww.TRACE.Println("Looking for the migration table")
		if !lib.HasMigrationTable(db) {
			jww.INFO.Println("No migration table found, exiting.")
			return
		}

		// Initialize 2 migrations arrays:
		// 1. The already applied migrations, as stated by the migrations table
		// 2. The available migrations, by looking into the directory given in
		//    the configuration file
		// Then, sort both by inverted migration version, to help the action loop
		// that come after.
		jww.TRACE.Println("Looking for applied migrations")
		applied, err := lib.GetAppliedMigrations(db)
		if err != nil {
			jww.ERROR.Println("Unable to get applied migrations:", err)
			return
		}

		jww.TRACE.Println("Looking for available migrations")
		available, err := lib.GetAvailableMigrations()
		if err != nil {
			jww.ERROR.Println("Unable to get available migrations:", err)
			return
		}

		sort.Sort(applied)
		sort.Sort(available)

		// Iterate over both array from the end until we have reached the begining
		// of both arrays. There is 3 cases:
		// 1. The applied version is greater than the available version, there is
		//    a missing migration. Report error then stop.
		// 2. The applied version is lesser than the available version, and :
		//    1. The applied version is the first (last) of the slice. Advance to
		//       the next available migration, as this is a not yet applied migration.
		//    2. The applied version isn't the first (last) of the slice. Report
		//       a missing migration then stop.
		// 3. Version both equals. Rollback the applied migration.
		// TODO Add an option to ignore out-of-order migrations
		jww.TRACE.Println("Analyzing migrations")
		var i, j int = len(applied) - 1, len(available) - 1
		for i >= 0 && j >= 0 {
			if applied[i].Version > available[j].Version {
				jww.ERROR.Println("Missing migration:", applied[i].File)
				return
			} else if applied[i].Version < available[j].Version {
				if i == len(applied)-1 {
					j--
				} else {
					jww.ERROR.Println("Found out-of-order new migration:", available[j].File)
					return
				}
			} else {
				// Rollback the migration.
				jww.INFO.Println("Reversing", applied[i].File)
				statements, err := applied[i].Scan("down")
				if err != nil {
					jww.ERROR.Println("Unable get down sections:", err)
					return
				}

				jww.TRACE.Println("Starting transaction")
				tx, err := db.Beginx()
				if err != nil {
					jww.ERROR.Println("Unable to start transaction:", err)
				}

				for i := len(statements) - 1; i >= 0; i-- {
					statement := statements[i]
					jww.INFO.Println("Executing statement:", statement)
					_, err := tx.Exec(statement)
					if err != nil {
						jww.ERROR.Println(err)
						jww.INFO.Println("Rollbacking")
						err := tx.Rollback()
						if err != nil {
							jww.ERROR.Println("Unable to rollback:", err)
							return
						}
						return
					}
				}

				jww.TRACE.Println("Removing entry in the migration table")
				_, err = tx.Exec("DELETE FROM migrations WHERE version = ?", applied[i].Version)
				if err != nil {
					jww.ERROR.Println("Unable to remove entry in the migrations table:", err)
					jww.INFO.Println("Rollbacking")
					err := tx.Rollback()
					if err != nil {
						jww.ERROR.Println("Unable to rollback:", err)
						return
					}
					return
				}

				jww.TRACE.Println("Committing")
				err = tx.Commit()
				if err != nil {
					jww.ERROR.Println("Unable to commit the transaction:", err)
					return
				}

				if viper.GetBool("reverse-all") == false {
					break
				}

				i--
				j--
			}
		}

		// Announce it's done
		jww.INFO.Println("Done")
	}),
}
