// Package aacautoupdate periodically grabs the latest version of the donor
// info XLS file, converts it into a CSV files, and grabs the needed values
// out of that file. After that, the info is placed into a JSON file that
// the web app reads from.
package main

import (
	"flag"
	"fmt"
	"syscall"

	"github.com/iAmSomeone2/aacautoupdate/update"
)

// Main sets up the main loop.
func main() {
	// Set up cmd line flags
	urlPtr := flag.String("source", "", "A web URL for accessing the patron data.")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Println("ERROR: A URL must be provided to use this program!")
		syscall.Exit(1)
	}

	fileName := update.CheckForUpdate(*urlPtr)
	fmt.Printf("Downloaded file located at: '%s'\n", fileName)

	//var cleanData string
	// If fileName is not empty, process the data in that file.
	if fileName != "" {
		// Continue work to process the data.
		//cleanData, _ = data.Clean(fileName)
	}

	// Wait for the next check.
}
