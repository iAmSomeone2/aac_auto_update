package update_test

import (
	"log"
	"os"
	"testing"

	"github.com/iAmSomeone2/aacautoupdate/update"
)

const (
	dlFileName string = "patrons_raw-html.txt"
)

func TestCheckForUpdate(t *testing.T) {
	fileName := update.CheckForUpdate(dlURL)

	if fileName != dlFileName {
		t.Error(
			"For", "TestCheckForUpdate()",
			"expected", dlFileName,
			"got", fileName,
		)
	}

	err := os.Remove(fileName)
	if err != nil {
		log.Panic(err)
	}
}
