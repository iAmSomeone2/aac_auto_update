// Package aacautoupdate periodically grabs the latest version of the donor
// info XLS file, converts it into a CSV files, and grabs the needed values
// out of that file. After that, the info is placed into a JSON file that
// the web app reads from.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/iAmSomeone2/aacautoupdate/data"
	"github.com/iAmSomeone2/aacautoupdate/update"
)

// Main sets up the main loop.
func main() {
	// Set up cmd line flags
	urlPtr := flag.String("source", "", "A web URL for accessing the patron data.")
	cleanPtr := flag.Bool("cleanrun", false, "Set this flag to clear the download cache.")

	flag.Parse()

	if *urlPtr == "" {
		log.Fatalln("ERROR: A URL must be provided to use this program!")
	}

	// If the cleanrun flag is set, delete the current and previous txt files
	if *cleanPtr {
		err := os.Remove("patrons_raw.txt")
		if err != nil {
			log.Println(err)
		}
		err = os.Remove("patrons_raw.old.txt")
		if err != nil {
			log.Println(err)
		}
	}

	fileName := update.CheckForUpdate(*urlPtr)
	fmt.Printf("Downloaded file located at: '%s'\n", fileName)

	//var cleanData string
	// If fileName is not empty, process the data in that file.
	if fileName != "" {
		// Continue work to process the data.
		cleanData, _ := data.Clean(fileName)
		patronList := data.NewPatronList(data.GetPatronData(cleanData))
		cellList := data.NewCellList(patronList)
		fmt.Println(patronList)
	}

	// Wait for the next check.
}
