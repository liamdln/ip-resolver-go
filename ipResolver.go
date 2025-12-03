package ipresolver

import (
	"bytes"
	"errors"
	"fmt"
	"maps"
	"net"
	"slices"
	"strings"

	"github.com/liamdln/ip-resolver-go/internal"
)

var ipData map[string]internal.Octet

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
func ResolveIp(ip string) (*internal.IPMetadata, error) {

	// ensure ip data has been loaded
	if len(ipData) < 1 {
		return nil, errors.New("no database was loaded, ensure you loaded a database first")
	}

	// ensure the ip is v4
	parsedIp := net.ParseIP(ip)
	if parsedIp.To4() == nil {
		return nil, errors.New("ip specified is not v4")
	}

	// split the IP into octets
	octets := strings.Split(ip, ".")
	finalOctets := ipData[octets[0]].Child[octets[1]].Child[octets[2]]
	fourthOctets := slices.Collect(maps.Keys(finalOctets.Child))

	// if there are no fourth octets, there is no IP that matches
	//
	// if there is only 1, then the IP falls into that range
	if len(fourthOctets) < 1 {
		return nil, fmt.Errorf("no range exists that contains %s", ip)
	} else if len(fourthOctets) == 1 {
		return finalOctets.Child[fourthOctets[0]].Value, nil
	}

	// if there is more than one entry in the fourth octet map, find the closest
	//
	// loop over the metadata and check if the ip is in the range
	for _, elem := range fourthOctets {
		// get the metadata from this fourth octet key
		// so if it's 10.0.0.128, look for finalOctets["128"]
		rangeStart := net.ParseIP(finalOctets.Child[elem].Value.RangeStart)
		rangeEnd := net.ParseIP(finalOctets.Child[elem].Value.RangeEnd)

		// check if the ip is in the range
		if bytes.Compare(parsedIp, rangeStart) >= 0 && bytes.Compare(parsedIp, rangeEnd) <= 0 {
			// ip is in range
			return finalOctets.Child[elem].Value, nil
		}
	}

	// ip address not found
	return nil, errors.New("ip address not found in database")

}
