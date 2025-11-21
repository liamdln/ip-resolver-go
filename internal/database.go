package internal

import (
	"encoding/csv"
	"os"
	"strconv"
)

type IPDatabase struct {
	RangeStart    string
	RangeEnd      string
	CountryCode   string
	ASNumber      int
	ASDescription string
}

// LoadDb loads a database file from a given path.
func LoadDb(filepath string) ([]IPDatabase, error) {

	dbFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	// close the db file once the func completes
	defer dbFile.Close()

	// setup csv reader with the deliemeter as a tab
	r := csv.NewReader(dbFile)
	r.Comma = '\t'

	// read the rows from the csv file
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	// create the ipdatabase array
	ipData := make([]IPDatabase, 0)

	// populate the ipdatabase array
	for _, elem := range rows {
		// try convert the as number to an int
		asNum, err := strconv.Atoi(elem[2])
		if err != nil {
			// skip this row
			continue
		}

		// add the data to the array
		ipData = append(ipData, IPDatabase{
			RangeStart:    elem[0],
			RangeEnd:      elem[1],
			ASNumber:      asNum,
			CountryCode:   elem[3],
			ASDescription: elem[4],
		})
	}

	return ipData, nil

}
