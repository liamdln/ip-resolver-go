package main

import (
	"fmt"
	"math/rand"
	"time"

	ipresolver "github.com/liamdln/ip-resolver-go"
)

func generateRandomIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256), rand.Intn(256))
}

func main() {

	// cloudflare IP
	ip := "241.161.115.14"

	// load the asn database
	err := ipresolver.LoadIPFile("./test.tsv")
	if err != nil {
		panic(err)
	}

	start := time.Now()
	// resolve the IP
	data, err := ipresolver.ResolveIp(ip)
	if err != nil {
		panic(err)
	}

	fmt.Println(data)

	fmt.Printf("took %s\n", time.Since(start))

	// for range 1000 {
	// 	ip := generateRandomIPv4()
	// 	fmt.Println("selected ip: ", ip)
	// 	data, err := ipresolver.ResolveIp(ip)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Println(data)
	// }

}
