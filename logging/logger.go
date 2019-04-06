package logging

import (
	"log"
	"os"
	"path"

	"github.com/coreos/go-systemd/journal"
)

// Logger is a struct used for simplifying the usage of the journalctl system.
// If the journactl system cannot be used, then a log file is output to a text
// file in /var/log.
type Logger struct {
	journalAvail bool
}

const logFile string = "/var/log/aac-auto-update/run.log"

// NewLogger retruns a pointer to an instance of the Logger struct. Use this
// instead of directly creating a Logger.
func NewLogger() *Logger {
	journalAvail := journal.Enabled()

	if !journalAvail {
		err := os.MkdirAll(path.Dir(logFile), 644)
		f, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}
		defer f.Close()

		log.SetOutput(f)
	}

	return &Logger{journalAvail: journalAvail}
}

// Println mimics the Println method from log.Println() or uses a journal
// equivalent.
func (logger *Logger) Println(v ...interface{}) error {
	if logger.journalAvail {
		return journal.Print(journal.PriInfo, "%v", v)
	}
	log.Println(v)
	return nil
}
