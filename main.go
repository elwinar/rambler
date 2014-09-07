package main

import (
	"fmt"
	"os"
)

// Dispatch the flags and run the appropriate command
func main() {
	err := command.Dispatch(os.Args[1:])
	if err != nil {
		fmt.Println(err)
	}

	return
}
