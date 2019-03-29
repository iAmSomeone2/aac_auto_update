package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// PatronList is a struct used for managing a list of Patrons.
type PatronList struct {
	patrons     []*Patron
	length      int
	totalRaised int
	totalCells  float32
}

// NewPatronList constructs a PatronList and returns a
// pointer to it. Only a slice of Patrons is required. All
// other values are computed from the list.
func NewPatronList(newPatrons []*Patron) *PatronList {
	var amtRaised int
	var cellNum float32
	patronNum := len(newPatrons)

	for _, patron := range newPatrons {
		amtRaised += patron.pledgeAmt
		cellNum += patron.cellAmt
	}

	return &PatronList{
		patrons:     newPatrons,
		length:      patronNum,
		totalRaised: amtRaised,
		totalCells:  cellNum,
	}
}

// MarshalJSON implements the MarshalJSON interface and allows for formatting
// the PatronList struct as JSON data.
func (patronList PatronList) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{\"patrons\":[")
	// Construct the portion of the data containing the patrons
	for i, patron := range patronList.patrons {
		patronJSON, err := patron.MarshalJSON()
		if err != nil {
			return nil, err
		}

		buffer.WriteString(string(patronJSON))

		// Add a comma between every patron except for the last one.
		if i < patronList.length-1 {
			buffer.WriteRune(',')
		}
	}
	buffer.WriteString("],")

	// Write in the data exclusive to the PatronList object
	// length field
	lenJSON, err := json.Marshal(patronList.length)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "length", string(lenJSON)))

	// total_raised field
	raisedJSON, err := json.Marshal(patronList.totalRaised)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "total_raised", string(raisedJSON)))

	// total_cells field
	cellsJSON, err := json.Marshal(patronList.totalCells)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "total_cells", string(cellsJSON)))

	buffer.WriteRune('}')
	return buffer.Bytes(), nil
}

// String returns a string version of the data PatronList represents
// MarshalJSON is used to construct the string.
func (patronList *PatronList) String() string {
	jsonStr, err := json.MarshalIndent(patronList, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	return string(jsonStr)
}
