package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

// Cells is a struct that represents all information regarding the cell array.
// This includes the PatronList, active cells, and their owners.
type CellList struct {
	patrons          *PatronList
	cells            []*Cell
	credit           float32
	remainingPatrons []int
}

// Cell is a struct representing an individual cell from the array. The id value
// is the int value referencing the associated cell in the array. The adoptee value
// is a slice of Patron pointers because it is possible for a cell to be owned by
// any number of patrons as long as their contributions equal out to the price per cell.
type Cell struct {
	id         int
	adopteeIDs []int
}

// NewCellList returns a pointer to a newly created CellList object. The PatronList is
// placed directly into the object. The Cell pointer slice is constructed based on the
// contents of the PatronList.
func NewCellList(list *PatronList) *CellList {
	// TODO: group adoptees that pay <$50 so that the cell splits make sense.

	// For each Patron in the PatronList, construct a Cell and determine which patrons are the adoptees.
	var adoptees []*Patron
	var runningCredit float32
	var cells []*Cell
	cellsIdx := 1

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
				cells = append(cells, &Cell{id: cellsIdx, adopteeIDs: []int{patron.id}})
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
		//adoptees = append(adoptees, patron)

		/*
			If runningCredit is at least 1, create and append a cell that is attributed to the accumulated patrons.
			This process is essentially the same as it is with an individual patron.
		*/
		cellNum = int(runningCredit)
		runningCredit = runningCredit - float32(cellNum)

		if cellNum > 0 {
			adoptees, remaining, credit := groupAdoptees(adoptees, runningCredit)
			var ids []int
			for _, patron := range adoptees {
				ids = append(ids, patron.id)
			}
			for i := 0; i < cellNum; i++ {
				cells = append(cells, &Cell{id: cellsIdx, adopteeIDs: ids})
				cellsIdx++
			}
			// Empty the adoptees slice
			adoptees = remaining
			runningCredit = credit
		}

		// If there is still a little credit left, attribute it to the current patron.
		if runningCredit > 0 {
			adoptees = append(adoptees, patron)
		}
	}

	var remainingIDs []int
	for _, patron := range adoptees {
		remainingIDs = append(remainingIDs, patron.id)
	}

	return &CellList{
		patrons:          list,
		cells:            cells,
		credit:           runningCredit,
		remainingPatrons: remainingIDs,
	}
}

// groupAdoptees groups the left over adoptees so that the amounts that they paid equal out to an even $50 as best as
// possible. The first *PatronList returned is the grouped adoptees, and the second is the remaining ones.
func groupAdoptees(adoptees []*Patron, credit float32) ([]*Patron, []*Patron, float32) {
	var groups []*Patron
	remaining := adoptees

	/*
		For each iteration of the loop, find the patron that has paid the most and determine how much is remaining to
		reach the cost of a single cell. Then find the next patron in the list that can fill what's left as close as
		possible. Continue until enough credit to count as a cell has been attributed. If getting an even 1 cell isn't
		possible, the roll over patron will have any additional credit over 1 set as their donation amount, and be
		put back into the 'remaining' list so that they can be used to fill in a gap later.
	*/

	for credit > 1 {
		// Find the top payer in 'Remaining'
		topPatron := remaining[0]
		var usedIdx []int
		var topIdx int
		for i, patron := range remaining {
			// Skip the first entry since we already have it
			if i == 0 {
				continue
			}

			// Compare each remaining entry against the current top
			if patron.cellAmt > topPatron.cellAmt {
				topPatron = patron
				topIdx = i
			}
		}
		usedIdx = append(usedIdx, topIdx)

		// Now, find the next top payers where their total cells are <= 1.
		owners := NewPatronList([]*Patron{topPatron})
		for owners.totalCells < 1 {
			topPatron = nil
			var topIdx int
			for i, patron := range remaining {
				// Skip any used entries since we already have them.
				skip := false
				for _, val := range usedIdx {
					if i == val {
						skip = true
						break
					}
				}

				if skip {
					continue
				}

				if topPatron == nil {
					topPatron = patron
					topIdx = i
				}

				// Compare each remaining entry against the current top
				if patron.cellAmt > topPatron.cellAmt {
					topPatron = patron
					topIdx = i
				}
			}
			// If we get here and topPatron is still nil, then we need to figure out if we have a full cell or not.
			if topPatron == nil {
				break
			}

			usedIdx = append(usedIdx, topIdx)
			owners.AddPatron(topPatron)
		}

	}

	return groups, remaining, credit
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

	// adoptee_ids field
	adopteesJSON, err := json.Marshal(cell.adopteeIDs)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "adoptee_ids", string(adopteesJSON)))

	buffer.WriteRune('}')
	//fmt.Println(string(buffer.Bytes()))
	return buffer.Bytes(), nil
}

func (list CellList) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")

	// Marshall in the Cell pointer slice
	buffer.WriteString("\"adopted_cells\":[")

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

	buffer.WriteString("],")
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
	buffer.WriteString(fmt.Sprintf("\"%s\":%s,", "remaining_patrons", string(remainingJSON)))

	// Marshal in the patron data
	//Marshall in the remainingPatrons PatronList
	patronsJSON, err := json.Marshal(list.patrons)
	if err != nil {
		return nil, err
	}
	buffer.WriteString(fmt.Sprintf("\"%s\":%s", "patron_list", string(patronsJSON)))

	buffer.WriteRune('}')
	// fmt.Println(string(buffer.Bytes()))
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

// ToJSONFile writes the contents of the CellList to a JSON-formatted text file.
func (list *CellList) ToJSONFile(fileName string) error {
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to the file.
	if err = ioutil.WriteFile(fileName, data, 0644); err != nil {
		return err
	}

	return nil
}
