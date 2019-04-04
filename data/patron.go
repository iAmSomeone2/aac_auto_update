// Package data provides functionality for retrieving relevant data from the
// exported data set. This data can be mapped to the Patron struct to allow for
// consistent data handling.
package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// Patron provides a structure for storing the values needed for updating the
// JSON file that the web app reads from.
type Patron struct {
	id         int
	pledgeTime time.Time
	anonymous  bool
	firstName  string
	lastName   string
	pledgeAmt  int
	cellAmt    float32
}

const (
	anonValIdx     int = 2
	fNameValIdx    int = 5
	timePledgedIdx int = 0
	pledgeValIdx   int = 32

	cellCost int = 50

	// Time used in this string must equate to Jan 2 15:04:05 MST 2006
	// Example time from patron data: 2019-03-31 08:21:16
	timeLayout string = "2006-01-02 15:04:05 MST"
	timeZone   string = "CST"
)

// NewPatron returns a new Patron struct based off of the values passed when the
// function is called. cellAmt is computed based on the pledge amount and may be
// any floating point value greater than 0.
func NewPatron(id int, pledgeTime string, anon bool, fName, lName string, pledgeAmt int) *Patron {
	cellNum := float32(pledgeAmt) / float32(cellCost)

	if anon {
		fName = "Anonymous"
		lName = "Donor"
	}

	// Create a time object from the imported time
	parsedTime, err := time.Parse(timeLayout, pledgeTime+" "+timeZone)
	if err != nil {
		log.Panic(err)
	}

	// fmt.Println(parsedTime)

	return &Patron{
		id:         id,
		pledgeTime: parsedTime,
		anonymous:  anon,
		firstName:  fName,
		lastName:   lName,
		pledgeAmt:  pledgeAmt,
		cellAmt:    cellNum,
	}
}

// MarshalJSON marshals the Patron struct into a JSON-compatible byte slice.
func (patron Patron) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	/*
		For each field in the Patron struct, construct a new string to be passed
		to the byte buffer.
	*/
	// id field
	idJSON, err := json.Marshal(patron.id)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "id", string(idJSON)))

	// pledgeTime field
	timeJSON, err := json.Marshal(patron.pledgeTime)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "pledge_time", string(timeJSON)))

	// anonymous field
	anonJSON, err := json.Marshal(patron.anonymous)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "anonymous", string(anonJSON)))

	// first_name field
	fNameJSON, err := json.Marshal(patron.firstName)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "first_name", string(fNameJSON)))

	// last_name field
	lNameJSON, err := json.Marshal(patron.lastName)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "last_name", string(lNameJSON)))

	// pledge_amt field
	pledgeJSON, err := json.Marshal(patron.pledgeAmt)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "pledge_amt", string(pledgeJSON)))

	// cell_num field
	cellJSON, err := json.Marshal(patron.cellAmt)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "cell_amt", string(cellJSON)))

	buffer.WriteString("}")
	return buffer.Bytes(), nil
}

// String returns the values contained in a Patron struct formatted so that
// it makes sense to read.
func (patron *Patron) String() string {
	//this := *patron

	jsonBytes, _ := json.Marshal(patron)

	return string(jsonBytes)
}
