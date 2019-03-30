package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
)

// Cells is a struct that represents all information regarding the cell array.
// This includes the PatronList, active cells, and their owners.
type CellList struct {
	patrons          *PatronList
	cells            []*Cell
	credit           float32
	remainingPatrons *PatronList
}

// Cell is a struct representing an individual cell from the array. The id value
// is the int value referencing the associated cell in the array. The adoptee value
// is a slice of Patron pointers because it is possible for a cell to be owned by
// any number of patrons as long as their contributions equal out to the price per cell.
type Cell struct {
	id      int
	adoptee []*Patron
}

// NewCellList returns a pointer to a newly created CellList object. The PatronList is
// placed directly into the object. The Cell pointer slice is constructed based on the
// contents of the PatronList.
func NewCellList(list *PatronList) *CellList {

	// For each Patron in the PatronList, construct a Cell and determine which patrons are the adoptees.
	var adoptees []*Patron
	var runningCredit float32
	var cells []*Cell
	var cellsIdx int

	for _, patron := range list.patrons {
		cellNum := int(patron.cellAmt) // Doing this conversion should chop off the decimal value (which we want to do)
		credit := patron.cellAmt - float32(cellNum)

		/*
			cellNum is only used for the current patron. The idea here is that cellNum will always be an int, and that
			by checking this first, we can assign the correct number of cells to an individual patron. This should even
			cover cases in which a patron pays enough for at least a number of cells but still has left over credit.
		*/
		if cellNum > 0 {
			// Create a Cell for each one that has been adopted
			for i := 0; i < cellNum; i++ {
				cells = append(cells, &Cell{id: cellsIdx, adoptee: []*Patron{patron}})
				cellsIdx++
			}
		}

		/*
			After all of the cells that are wholly owned by a patron are created and added to the cells slice, if the
			credit remaining is between 0 and 1 cell it is added onto the runningCredit pile, and this patron is added
			onto the running list of adoptees.
		*/
		if credit == 0 {
			continue // No credit remains, so just skip to the next patron.
		}

		// Append the appropriate values into the running totals.
		runningCredit += credit
		adoptees = append(adoptees, patron)

		/*
			If runningCredit is at least 1, create and append a cell that is attributed to the accumulated patrons.
			This process is essentially the same as it is with an individual patron.
		*/
		cellNum = int(runningCredit)
		runningCredit = runningCredit - float32(cellNum)

		if cellNum > 0 {
			for i := 0; i < cellNum; i++ {
				cells = append(cells, &Cell{id: cellsIdx, adoptee: adoptees})
				cellsIdx++
			}
			// Empty the adoptees slice
			adoptees = nil
		}

		// If there is still a little credit left, attribute it to the current patron.
		if runningCredit > 0 {
			adoptees = append(adoptees, patron)
		}
	}

	return &CellList{
		patrons:          list,
		cells:            cells,
		credit:           runningCredit,
		remainingPatrons: NewPatronList(adoptees),
	}
}

// MarshallJSON formats the contents of the Cell struct so that it may be used in JSON data.
func (cell Cell) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	// total_cells field
	idJSON, err := json.Marshal(cell.id)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "id", string(idJSON)))

	buffer.WriteString("adoptee:[")
	for i, adoptee := range cell.adoptee {
		adoptJSON, err := json.Marshal(adoptee)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(string(adoptJSON))

		// Add a comma between every adoptee except for the last one.
		if i < len(cell.adoptee)-1 {
			buffer.WriteRune(',')
		}
	}

	buffer.WriteString("]}")
	return buffer.Bytes(), nil
}

func (list CellList) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	// Marshall in the main PatronList object
	patronsJSON, err := json.Marshal(list.patrons)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "patrons", string(patronsJSON)))

	// Marshall in the Cell pointer slice
	buffer.WriteString("adopted_cells:[")

	for i, cell := range list.cells {
		cellJSON, err := json.Marshal(cell)
		if err != nil {
			return nil, err
		}
		buffer.WriteString(string(cellJSON))

		if i < len(list.cells)-1 {
			buffer.WriteRune(',')
		}
	}

	buffer.WriteRune(']')
	// Marshall in the credit value
	creditJSON, err := json.Marshal(list.credit)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "credit", string(creditJSON)))

	//Marshall in the remainingPatrons PatronList
	remainingJSON, err := json.Marshal(list.remainingPatrons)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "remaining_patrons", string(remainingJSON)))

	buffer.WriteRune('}')
	return buffer.Bytes(), nil
}

// String returns a stringified version of the MarshalJSON output of Cell
func (cell Cell) String() string {
	out, err := json.MarshalIndent(cell, "", "  ")
	if err != nil {
		log.Panic(err)
	}

	return string(out)
}

// String returns a stringified version of the MarshalJSON output of CellList
func (list CellList) String() string {
	out, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		log.Panic(err)
	}

	return string(out)
}
