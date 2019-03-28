package data

import (
	"io/ioutil"
	"strings"
)

var replacementPairs = []string{
	"\"\"", " ", // Double quotes are replaced with a single space.
	"var data = ", "", // "var data = " is removed.
	"\"", "", // Individual quotation marks are removed.
	";", "", // Trailing semicolon removed.
}

// Clean reads the data from the file which the fileName argument is pointing to
// and places it into a string for initial processing. Extraneous characters and
// formatting is removed before the string is returned.
func Clean(fileName string) (string, error) {
	var cleanStr string
	// Read the file into memory and and assign it's data to a string for processing.
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return cleanStr, err
	}
	content := string(file)

	// Process the string using a separate function
	cleanStr = formatForUsage(content)

	return cleanStr, nil
}

// formatForUsage uses the strings library to trim all extraneous characters
// from the input string so that it can be effectively split into a string slice.
func formatForUsage(content string) string {
	var result string
	// First remove extraneous spaces
	result = strings.TrimSpace(content)

	// Next, build a replacer that will format our string
	replacer := strings.NewReplacer(replacementPairs...)
	result = replacer.Replace(result)

	return result
}

// GetPatronData takes in a string and returns a slice of
// Patron structs.
func GetPatronData(rawData string) []*Patron {
	// First, split the data into a 1D slice of strings using
}
