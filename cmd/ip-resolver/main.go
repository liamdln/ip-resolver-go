package main

import (
	"fmt"

	ipresolve "github.com/liamdln/ip-resolver-go"
)

func main() {

	// cloudflare IP
	ip := "1.1.1.1"

	// load the asn database
	err := ipresolve.LoadIPFile("./test.tsv")
	if err != nil {
		panic(err)
	}

	// resolve the IP
	data, err := ipresolve.ResolveIp(ip)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

}
