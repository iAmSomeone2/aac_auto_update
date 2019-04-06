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
		err := os.MkdirAll(path.Dir(logFile), os.ModeDir|os.ModePerm)
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
// equivalent, if available.
func (logger *Logger) Println(v ...interface{}) {
	logger.Printf("%v\n", v...)
}

// Printf mimics the Printf method from log.Printf() or uses a journal
// equivalent, if available.
func (logger *Logger) Printf(format string, v ...interface{}) {
	if logger.journalAvail {
		if err := journal.Print(journal.PriInfo, format, v...); err != nil {
			panic(err)
		}
	} else {
		log.Printf(format, v...)
	}
}

// Panic logs the interface and calls panic()
func (logger *Logger) Panic(v ...interface{}) {
	logger.Panicf("%v", v...)
}

// Panicf mimics the Panicf method from log.Panicf() or uses a journal
// equivalent, if available.
func (logger *Logger) Panicf(format string, v ...interface{}) {
	if logger.journalAvail {
		if err := journal.Print(journal.PriErr, format, v...); err != nil {
			panic(err)
		}
		panic(v)
	}
	log.Panicf(format, v...)
}

// Fatal mimics the Fatal method from log or uses a journal equivalent,
// if available.
func (logger *Logger) Fatal(v ...interface{}) {
	logger.Fatalf("%v", v...)
}

// Fatalln mimics the Fatalln method from log or uses a journal equivalent,
// if available.
func (logger *Logger) Fatalln(v ...interface{}) {
	logger.Fatalf("%v\n", v...)
}

// Fatalf mimics the Fatalf method from log.Fatalf or uses a journal
// equivalent, if available.
func (logger *Logger) Fatalf(format string, v ...interface{}) {
	if logger.journalAvail {
		if err := journal.Print(journal.PriCrit, format, v...); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
	log.Fatalf(format, v...)
}

// Warnf will log a warning to the system journal if it's available. If
// the system jounal is not available, then a warning is printed to stderr.
func (logger *Logger) Warnf(format string, v ...interface{}) {
	if logger.journalAvail {
		if err := journal.Print(journal.PriWarning, format, v...); err != nil {
			panic(err)
		}
	} else {
		log.Printf("WARN: %v\n", v...)
	}
}

// Warnln formats the warning string as a singular line and sends it to logger.Warnf()
func (logger *Logger) Warnln(v ...interface{}) {
	logger.Warnf("%v\n", v...)
}
