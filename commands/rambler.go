package commands

import (
	"github.com/spf13/cobra"
	"rambler/lib"
)

var (
	Rambler = &cobra.Command{
		Use:   "rambler",
		Short: "Rambler is a simple and language-independant SQL schema migration tool",
	}
)

func init() {
	// Add the common flags to the command persistent flags, so they are shared
	// between all subcommands
	*Rambler.PersistentFlags() = lib.Flags
	
}
