package serve

import (
	"io/ioutil"
	"net/http"

	"github.com/iAmSomeone2/aacautoupdate/logging"
)

const (
	dataLoc string = "/var/www/cell.bdavidson.dev/html/data/data.json"
)

func servePatronData(w http.ResponseWriter, r *http.Request) {
	/*
		Since we just need to send the raw JSON data, we should be able to
		read in the file, and serve the byte stream.
	*/

	data, err := ioutil.ReadFile(dataLoc)
	if err != nil {
		logger := logging.NewLogger()
		logger.Warnf("%v", err)
	}
	w.Write(data)
}

// StartServer sets up the simple HTTP server for handling data requests.
func StartServer() {
	logger := logging.NewLogger()
	logger.Printf("Data server started on separate thread.\n")
	http.HandleFunc("/patron-data", servePatronData)
	// ListenAndServe should be changed to the TLS variant for prod.
	if err := http.ListenAndServe(":8080", nil); err != nil {
		logger.Warnf("%v", err)
	}
}
