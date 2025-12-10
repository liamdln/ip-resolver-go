package ipresolver

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/liamdln/ip-resolver-go/internal"
)

var ipData map[string]internal.Octet

// the max number of times a search for a valid octet can loop for
const MAX_OCTET_SEARCH_ITERATIONS = 256

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

	// define a variable for errors
	var err error

	// check if the octet exists in the table.
	//
	// it may not as some IP ranges take multiple steps over an ocetet, for example,
	// range start: 1.0.0.0, range end: 3.0.0.0 (2.0.0.0 will not exist in the map)
	//
	// we know that if we subtract (staying above 0), we will get to the range start.
	//
	// check the first octet

	// TODO:
	// change this
	// so first check if the octet with the exact IP exists,
	// if not, find the min closest octet array to the ip octet and check its children
	// if the next octet is in the array (or if it is in the range of min octet - max octet), reloop
	// and set the searching octet to the next array (or the octet array of the min closest ip octet).
	// keep doing this 4 times then find a value.
	firstOctet := ipData[octets[0]]
	if firstOctet.Child == nil {
		firstOctet, err = findOctetRangeStart(octets[0], ipData)
		if err != nil {
			return nil, err
		}
	}

	// second octet
	secondOctet := firstOctet.Child[octets[1]]
	if secondOctet.Child == nil {
		secondOctet, err = findOctetRangeStart(octets[1], firstOctet.Child)
		if err != nil {
			return nil, err
		}
	}

	// third octet
	thirdOctet := secondOctet.Child[octets[2]]
	if thirdOctet.Child == nil {
		thirdOctet, err = findOctetRangeStart(octets[2], secondOctet.Child)
		if err != nil {
			return nil, err
		}
	}

	// fourth octet
	fourthOctet := thirdOctet.Child[octets[3]]
	if fourthOctet.Child == nil {
		fourthOctet, err = findOctetRangeStart(octets[3], thirdOctet.Child)
		if err != nil {
			return nil, err
		}
	}

	if !octetValueIsValid(fourthOctet.Value) {
		return nil, fmt.Errorf("ip %s not in table", ip)
	}

	return &fourthOctet.Value, nil

}

func findOctetRangeStart(octet string, children map[string]internal.Octet) (internal.Octet, error) {
	// convert the octet to a number
	octetInt, err := strconv.Atoi(octet)
	if err != nil {
		return internal.Octet{}, fmt.Errorf("could not convert octet %s to int: %s", octet, err.Error())
	}

	var newOctetStr string

	for range MAX_OCTET_SEARCH_ITERATIONS {
		newOctetStr = fmt.Sprintf("%d", octetInt)
		if octetInt < 0 {
			return internal.Octet{}, errors.New("octet searched reached 0, cannot search lower than 0")
		}
		if children[newOctetStr].Child != nil || octetValueIsValid(children[newOctetStr].Value) {
			break
		}
		octetInt -= 1
	}

	return children[newOctetStr], nil

}

func octetValueIsValid(metadata internal.IPMetadata) bool {
	return metadata.RangeStart != "" && metadata.RangeEnd != "" && metadata.CountryCode != "" && metadata.ASDescription != ""
}
