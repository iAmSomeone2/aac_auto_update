// Package data provides functionality for retrieving relevant data from the
// exported data set. This data can be mapped to the Patron struct to allow for
// consistent data handling.
package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

// Patron provides a structure for storing the values needed for updating the
// JSON file that the web app reads from.
type Patron struct {
	anonymous bool
	firstName string
	lastName  string
	pledgeAmt int
	cellNum   float32
}

// PatronList aliases a slice of Patron references.
type PatronList []*Patron

const (
	anonValIdx  int = 2
	fNameValIdx int = 5
	//lNameValIdx  int = 7
	pledgeValIdx int = 32

	cellCost int = 50
)

// NewPatron returns a new Patron struct based off of the values passed when the
// function is called. cellAmt is computed based on the pledge amount and may be
// any floating point value greater than 0.
func NewPatron(anon bool, fName, lName string, pledgeAmt int) *Patron {
	cellNum := float32(pledgeAmt) / float32(cellCost)

	return &Patron{
		anonymous: anon,
		firstName: fName,
		lastName:  lName,
		pledgeAmt: pledgeAmt,
		cellNum:   cellNum,
	}
}

/*
	TODO: Use MarshalJSON and UnmarshalJSON interfaces to make Patron compatible
	with JSON encoding and decoding.
*/

// MarshalJSON marshals the Patron struct into a JSON-compatible byte slice
func (patron *Patron) MarshalJSON() ([]byte, error) {
	this := *patron
	buffer := bytes.NewBufferString("{")

	/*
		For each field in the Patron struct, construct a new string to be passed
		to the byte buffer.
	*/
	anonJSON, err := json.Marshal(this.anonymous)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "anonymous", string(anonJSON)))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// String returns the values contained in a Patron struct formatted so that
// it makes sense to read.
func (patron *Patron) String() string {
	this := *patron
	var strBuilder strings.Builder

	fmt.Fprintf(&strBuilder,
		"{Patron: %s %s Pledge: $%d Cells Adopted: %f Anonymous: %t}",
		this.firstName,
		this.lastName,
		this.pledgeAmt,
		this.cellNum,
		this.anonymous,
	)

	return strBuilder.String()
}
