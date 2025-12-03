package internal

import (
	"encoding/csv"
	"os"
	"strconv"
	"strings"
)

type IPMetadata struct {
	RangeStart    string
	RangeEnd      string
	CountryCode   string
	ASNumber      int
	ASDescription string
}

// Octet is the struct that contains either IPMetadata or Child map
type Octet struct {
	Child map[string]*Octet
	Value *IPMetadata
}

// LoadDb loads a database file from a given path.
func LoadDb(filepath string) (map[string]Octet, error) {

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

	ips := make(map[string]Octet)

	// loop over each row and store the metadata in the ip map
	for _, elem := range rows {

		// elem[0] = range start
		// elem[1] = range end
		// elem[2] = ASNumber
		// elem[3] = country code
		// elem[5] = ASDescription

		// split the ip into octets
		octets := strings.Split(elem[0], ".")

		// try convert the as number to an int
		asNum, err := strconv.Atoi(elem[2])
		if err != nil {
			// skip this row
			continue
		}

		// ensure the map contains the first octet, if not, initialise a map
		if _, ok := ips[octets[0]]; !ok {
			ips[octets[0]] = Octet{
				Child: make(map[string]*Octet),
			}
		}

		// ensure the map contains the second octet, if not, initialise a map
		if _, ok := ips[octets[0]].Child[octets[1]]; !ok {
			ips[octets[0]].Child[octets[1]] = &Octet{
				Child: make(map[string]*Octet),
			}
		}

		// ensure the map contains the third octet, if not, initialise a map
		if _, ok := ips[octets[0]].Child[octets[1]].Child[octets[2]]; !ok {
			ips[octets[0]].Child[octets[1]].Child[octets[2]] = &Octet{
				Child: make(map[string]*Octet),
			}
		}

		// put the fourth octet in, as we know it won't exist for this octet pattern yet
		ips[octets[0]].Child[octets[1]].Child[octets[2]].Child[octets[3]] = &Octet{
			Value: &IPMetadata{
				RangeStart:    elem[0],
				RangeEnd:      elem[1],
				ASNumber:      asNum,
				CountryCode:   elem[3],
				ASDescription: elem[4],
			},
		}
	}

	return ips, nil

}
