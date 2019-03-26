// Package data provides functionality for retrieving relevant data from the
// exported data set. This data can be mapped to the Patron struct to allow for
// consistent data handling.
package data

// Patron provides a structure for storing the values needed for updating the
// JSON file that the web app reads from.
type Patron struct {
	anonymous bool
	firstName string
	lastName  string
	pledgeAmt int
	cellAmt   int
}

const (
	annonValIdx  int = 2
	fNameValIdx  int = 5
	lNameValIdx  int = 7
	pledgeValIdx int = 30
)
