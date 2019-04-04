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
	remainingPatrons map[int]*Patron
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
	creditPatrons := make(map[int]*Patron)
	var credit float32
	var cells []*Cell
	cellsIdx := 1

	for _, patron := range list.patrons {
		// This conversion will chop off any decimal values.
		for i := 0; i < int(patron.cellAmt); i++ {
			cells = append(cells, &Cell{id: cellsIdx, adopteeIDs: []int{patron.id}})
			cellsIdx++
		}

		// Throw any patron that hasn't paid enough for a cell into the creditIDs list
		if patron.cellAmt < 1 {
			creditPatrons[patron.id] = patron
			credit += patron.cellAmt
		}

		// If credit is >= 1 we should try to pair the extra patrons.
		if credit >= 1 {
			groups, remaining := groupPatrons(creditPatrons)
			// Create any new cells from the resulting groups
			if _, hasZero := groups[0]; !hasZero {
				for _, group := range groups {
					fmt.Println(groups)
					cells = append(cells, &Cell{id: cellsIdx, adopteeIDs: group})
					cellsIdx++
				}
				creditPatrons = remaining // Go won't let me assign this at the function call for some reason.

				// Here we update the remaining credit
				var newCredit float32
				for _, creditPatron := range creditPatrons {
					newCredit += creditPatron.cellAmt
				}
				credit = newCredit
			}
		}
	}

	return &CellList{
		patrons:          list,
		cells:            cells,
		credit:           credit,
		remainingPatrons: creditPatrons,
	}
}

// groupPatrons takes in a list of Patrons and groups them into a map if they can be put together to equal the
// value of a single cell. Any leftover Patrons will have their IDs returned separately.
func groupPatrons(patronMap map[int]*Patron) (map[int][]int, map[int]*Patron) {
	groups := make(map[int][]int)
	remains := make(map[int]*Patron)
	paired := make(map[*Patron]bool)

	count := 0
	for _, iPatron := range patronMap {
		// First, we should pair up anyone who has donated half of the value of a cell.
		if !paired[iPatron] {
			var group []int
			if iPatron.cellAmt == 0.5 {
				group = append(group, iPatron.id)
				//Check the remainder of the patronMap
				for _, jPatron := range patronMap {
					if !paired[jPatron] && jPatron != iPatron {
						if jPatron.cellAmt == 0.5 {
							group = append(group, jPatron.id)
							paired[iPatron] = true
							paired[jPatron] = true
							count++
							break
						}
					}
				}
			}
			groups[count] = group
		}
	}

	// Put the ID of every patron that wasn't paired into the "remains" array.
	for id, patron := range patronMap {
		if !paired[patron] {
			remains[id] = patron
		}
	}

	return groups, remains
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
	// TODO: Update this so that we just grab the Patron IDs.
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
