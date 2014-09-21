package commands

import (
	"github.com/elwinar/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"rambler/lib"
	"sort"
)

func init() {
	// Add the command to the main command
	Rambler.AddCommand(Apply)
	
	// Set the default configuration
// 	viper.SetDefault("all", false)
	
	// Add ubiquitous flags to the main command
// 	Apply.Flags().BoolP("all", "a", false, "apply all migrations")
	
	// Set overrides from the command-line to viper
// 	override("all", Apply.Flags().Lookup("all"))
}

var Apply = &cobra.Command{
	Use:   "apply",
	Short: "Apply one migration",
	Run:   do(func (cmd *cobra.Command, args []string) {
		// Open the database connection
		jww.TRACE.Println("Openning database connection")
		db, err := lib.GetDB()
		defer db.Close()
		if err != nil {
			jww.ERROR.Println("Unable to open database connection:", err)
			return
		}
		
		// Look for the migration table, and create it if not found
		jww.TRACE.Println("Looking for the migration table")
		if !lib.HasMigrationTable(db) {
			jww.INFO.Println("Migration table not found, creating it")
			err := lib.CreateMigrationTable(db)
			if err != nil {
				jww.ERROR.Println("Unable to create migration table:", err)
				return
			}
		}
		
		// Get the migrations already applied
		jww.TRACE.Println("Looking for applied migrations")
		applied, err := lib.GetAppliedMigrations(db)
		if err != nil {
			jww.ERROR.Println("Unable to get applied migrations:", err)
			return
		}
		sort.Sort(applied)
		
		// Get the migrations available
		jww.TRACE.Println("Looking for available migrations")
		available, err := lib.GetAvailableMigrations()
		if err != nil {
			jww.ERROR.Println("Unable to get available migrations:", err)
			return
		}
		sort.Sort(available)
		
		// Iterate over the applied migrations to detect out-of-order and missing migrations
		jww.TRACE.Println("Analyzing migrations")
		var i int
		for i = 0; i < len(applied) && i < len(available); i++ {
			if applied[i].Version > available[i].Version {
				jww.ERROR.Println("Found out-of-order new migration:", available[i].File)
				return
			} else if applied[i].Version < available[i].Version {
				jww.ERROR.Println("Missing migration:", applied[i].File)
				return
			} else {
				jww.TRACE.Println("Found already applied migration:", available[i].File)
			}
		}
		
		// If we are at both ends, there is nothing to do
		if i == len(available) && i == len(applied) {
			jww.INFO.Println("Nothing to do")
			return
		}
		
		// If there is remaining applied migrations but no available migrations, we are missing some
		if i < len(applied) {
			jww.ERROR.Println("Missing migration:", applied[i].File)
			return
		}
		
		// If there is remaining available migrations, get the statements of the first
		jww.INFO.Println("Applying", available[i].File)
		statements, err := available[i].Scan("up")
		if err != nil {
			jww.ERROR.Println("Unable get up sections:", err)
			return
		}
		
		// Start a transaction
		jww.TRACE.Println("Starting transaction")
		tx, err := db.Beginx()
		if err != nil {
			jww.ERROR.Println("Unable to start transaction:", err)
		}
		
		// Loop over the statements and execute them
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
		
		// Insert the migration in the table
		jww.TRACE.Println("Adding entry in the migration table")
		_, err = tx.Exec("INSERT INTO migrations (version, date, description, file) VALUES (?,?,?,?)", available[i].Version, available[i].Date.Format("2006-01-02 15:04:05"), available[i].Description, available[i].File)
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
		
		// Commit the transaction
		jww.TRACE.Println("Commiting")
		err = tx.Commit()
		if err != nil {
			jww.ERROR.Println("Unable to commit the transaction:", err)
			return
		}
		
		// Announce it's done
		jww.INFO.Println("Done")
	}),
}
