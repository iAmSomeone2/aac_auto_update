// Package aacautoupdate periodically grabs the latest version of the donor
// info XLS file, converts it into a CSV files, and grabs the needed values
// out of that file. After that, the info is placed into a JSON file that
// the web app reads from.
package main

import (
	"flag"
	"log"
	"os"
	"path"

	//"github.com/iAmSomeone2/aacautoupdate/data"
	"./data"
	"./update"
)

const (
	outputFile string = "data.json"
)

// Main sets up the main loop.
func main() {
	// Set up cmd line flags
	urlPtr := flag.String("source", "", "A web URL for accessing the patron data.")
	cleanPtr := flag.Bool("cleanrun", false, "Set this flag to clear the download cache.")
	outPtr := flag.String("out", "./", "The directory in which to place the data.json file.")

	flag.Parse()

	if *urlPtr == "" {
		log.Fatalln("ERROR: A URL must be provided to use this program!")
	}

	// If the cleanrun flag is set, delete the current and previous txt files
	if *cleanPtr {
		cacheDir := update.GetCacheDir()
		cacheDir = path.Join(cacheDir, update.AppDir)
		err := os.Remove(path.Join(cacheDir, update.BaseFileName))
		if err != nil {
			log.Println(err)
		}
		err = os.Remove(path.Join(cacheDir, update.OldFileName))
		if err != nil {
			log.Println(err)
		}
	}

	outputPath := path.Join(*outPtr, outputFile)

	fileName := update.CheckForUpdate(*urlPtr)

	//var cleanData string
	// If fileName is not empty, process the data in that file.
	if fileName != "" {
		log.Printf("Downloaded file located at: '%s'\n", fileName)
		// Continue work to process the data.
		cleanData, _ := data.Clean(fileName)
		patronList := data.NewPatronList(data.GetPatronData(cleanData))
		cellList := data.NewCellList(patronList)
		if err := cellList.ToJSONFile(outputPath); err != nil {
			log.Panic(err)
		} else {
			log.Printf("Data written to %s\n", outputFile)
		}
	} else {
		log.Printf("Nothing to do. Will check again soon.\n")
	}

	// Wait for the next check.
}
