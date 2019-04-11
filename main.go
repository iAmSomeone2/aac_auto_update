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
	"time"

	"github.com/iAmSomeone2/aacautoupdate/data"
	"github.com/iAmSomeone2/aacautoupdate/logging"
	"github.com/iAmSomeone2/aacautoupdate/serve"
	"github.com/iAmSomeone2/aacautoupdate/update"
)

const (
	outputFile string = "data.json"
	defaultURL string = "https://campaigns.communityfunded.com/download-supporters/?p_id=26458"
	defaultDir string = "/var/www/cell.bdavidson.dev/html/data"
)

// Main sets up the main loop.
func main() {
	// Set up cmd line flags
	urlPtr := flag.String("source", defaultURL, "A web URL for accessing the patron data.")
	cleanPtr := flag.Bool("cleanrun", false, "Set this flag to clear the download cache.")
	outPtr := flag.String("out", defaultDir, "The directory in which to place the data.json file.")
	waitPtr := flag.Int64("wait", 5, "An integer value representing the number of minutes to wait between checks.")

	flag.Parse()

	if *urlPtr == "" {
		log.Fatalln("ERROR: A URL must be provided to use this program!")
	}

	logger := logging.NewLogger()

	// If the cleanrun flag is set, delete the current and previous txt files
	if *cleanPtr {
		cacheDir := update.GetCacheDir()
		cacheDir = path.Join(cacheDir, update.AppDir)
		err := os.Remove(path.Join(cacheDir, update.BaseFileName))
		if err != nil {
			logger.Warnln(err)
		}
		err = os.Remove(path.Join(cacheDir, update.OldFileName))
		if err != nil {
			logger.Warnln(err)
		}
	}

	outputPath := path.Join(*outPtr, outputFile)

	// Start HTTP server on a separate thread to serve the data file.
	go serve.StartServer()

	startLoop := true
	waitTime := time.Duration(*waitPtr * int64(time.Minute))
	for { // Run through this every five minutes.
		updateTimer := time.NewTimer(waitTime)

		timerStop := false
		if startLoop {
			logger.Println("Immediately pulling update for initial run.")
			timerStop = updateTimer.Stop()
			startLoop = false
		}
		if !timerStop {
			<-updateTimer.C
		}
		fileName := update.CheckForUpdate(*urlPtr)

		// If fileName is not empty, process the data in that file.
		if fileName != "" {
			log.Printf("Downloaded file located at: '%s'\n", fileName)
			// Continue work to process the data.
			cleanData, _ := data.Clean(fileName)
			patronList := data.NewPatronList(data.GetPatronData(cleanData))
			cellList := data.NewCellList(patronList)
			if err := cellList.ToJSONFile(outputPath); err != nil {
				logger.Fatal(err)
			} else {
				logger.Printf("Data written to %s\n", outputFile)
			}
		} else {
			logger.Printf("Nothing to do. Will check again soon.\n")
		}

		// Wait for the next check.
		var s string
		if *waitPtr != 1 {
			s = "s"
		} else {
			s = ""
		}
		logger.Printf("Check finished. Waiting %d minute%s...\n", *waitPtr, s)
	}
}
