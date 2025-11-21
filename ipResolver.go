package ipresolver

import (
	"bytes"
	"errors"
	"net"

	"github.com/liamdln/ip-resolve-go/internal"
)

var ipData []internal.IPDatabase

// LoadIPFile loads an IP-ASN file
func LoadIPFile(fileLocation string) error {

	// load the ip data
	var err error
	ipData, err = internal.LoadDb(fileLocation)
	if err != nil {
		return err
	}

	return nil

}

// ResolveIp resolves a given IP
func ResolveIp(ip string) (*internal.IPDatabase, error) {

	// ensure ip data has been loaded
	if len(ipData) < 1 {
		return nil, errors.New("no database was loaded, ensure you loaded a database first")
	}

	// parse the IP and ensure it's v4
	parsedIp := net.ParseIP(ip)
	if parsedIp.To4() == nil {
		return nil, errors.New("ip specified is not v4")
	}

	// loop over the ips to find a range the speicifed IP falls into
	for _, elem := range ipData {
		rangeStart := net.ParseIP(elem.RangeStart)
		rangeEnd := net.ParseIP(elem.RangeEnd)

		// check if the ip is in the range
		if bytes.Compare(parsedIp, rangeStart) >= 0 && bytes.Compare(parsedIp, rangeEnd) <= 0 {
			// ip is in range
			return &elem, nil
		}
	}

	// ip address not found
	return nil, errors.New("ip address not found in database")

}
