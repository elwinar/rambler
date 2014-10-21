package commands

import (
	"fmt"
	"github.com/elwinar/cobra"
	"github.com/elwinar/rambler/lib"
	"github.com/elwinar/viper"
	jww "github.com/spf13/jwalterweatherman"
	"sort"
)

func init() {
	// Add the command to the main command
	Rambler.AddCommand(Apply)

	// Set the default configuration
	viper.SetDefault("all", false)

	// Add ubiquitous flags to the main command
	Apply.Flags().BoolP("all", "a", false, "apply all migrations")

	// Set overrides from the command-line to viper
	override("apply-all", Apply.Flags().Lookup("all"))
}

var Apply = &cobra.Command{
	Use:   "apply",
	Short: "Apply the next migration",
	Run: do(func(cmd *cobra.Command, args []string) {
		// Start by opening the database connection and looking for the migration
		// table. If not found, we have to create it.
		// TODO Add an option to create/not create the table (default to be determined)
		jww.TRACE.Println("Openning database connection")
		db, err := lib.GetDB()
		if err != nil {
			jww.ERROR.Println("Unable to open database connection:", err)
			return
		}
		defer db.Close()

		jww.TRACE.Println("Looking for the migration table")
		if !lib.HasMigrationTable(db) {
			jww.INFO.Println("Migration table not found, creating it")
			err := lib.CreateMigrationTable(db)
			if err != nil {
				jww.ERROR.Println("Unable to create migration table:", err)
				return
			}
		}

		// Initialize 2 migrations arrays:
		// 1. The already applied migrations, as stated by the migrations table
		// 2. The available migrations, by looking into the directory given in
		//    the configuration file
		// Then, sort both by migration version, to help the action loop that
		// come after.
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

		// Iterate over both array until we have reached the end of both arrays.
		// There is 5 cases:
		// 1. After the end of applied, but not of available, this is a new
		//    migration in available. Apply it.
		// 2. After the end of available but not of applied, this is a missing
		//    migration. Report error then stop.
		// 3. Applied version greater than available version, the available migration
		//    is probably new, but added at the wrong place. Report error then stop.
		// 4. Applied version lesser than available version, the applied migration
		//    is missing in the available array. Report error then stop.
		// 5. Both versions are equal, the applied version is found, nothing
		//    to do. Log in the trace and continue.
		// TODO Add an option to apply out-of-order migrations
		jww.TRACE.Println("Analyzing migrations")
		var i int
		for i = 0; i < len(applied) || i < len(available); i++ {
			if i >= len(applied) && i < len(available) {
				// Apply the migration.
				jww.INFO.Println("Applying", available[i].File)
				statements, err := available[i].Scan("up")
				if err != nil {
					jww.ERROR.Println("Unable to get up sections:", err)
					return
				}

				jww.TRACE.Println("Starting transaction")
				tx, err := db.Beginx()
				if err != nil {
					jww.ERROR.Println("Unable to start transaction:", err)
				}

				for _, statement := range statements {
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

				jww.TRACE.Println("Adding entry in the migration table")
				sqlStr := fmt.Sprintf("INSERT INTO migrations (version, date, description, file) VALUES ('%d','%s','%s','%s')", available[i].Version, available[i].Date.Format("2006-01-02 15:04:05"), available[i].Description, available[i].File)
				_, err = tx.Exec(sqlStr)
				if err != nil {
					jww.ERROR.Println("Unable to write entry in the migrations table:", err)
					jww.INFO.Println("Rollbacking")
					err := tx.Rollback()
					if err != nil {
						jww.ERROR.Println("Unable to rollback:", err)
						return
					}
					return
				}

				jww.TRACE.Println("Commiting")
				err = tx.Commit()
				if err != nil {
					jww.ERROR.Println("Unable to commit the transaction:", err)
					return
				}

				if viper.GetBool("apply-all") == false {
					break
				}
			} else if i >= len(available) && i < len(applied) {
				jww.ERROR.Println("Missing migration:", applied[i].File)
				return
			} else if applied[i].Version > available[i].Version {
				jww.ERROR.Println("Found out-of-order new migration:", available[i].File)
				return
			} else if applied[i].Version < available[i].Version {
				jww.ERROR.Println("Missing migration:", applied[i].File)
				return
			} else {
				jww.TRACE.Println("Found already applied migration:", available[i].File)
			}
		}

		// Announce it's done
		jww.INFO.Println("Done")
	}),
}
