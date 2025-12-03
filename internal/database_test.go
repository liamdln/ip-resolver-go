package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ensure the database file loads correctly
func TestLoadDb(t *testing.T) {

	// create a tab separated row
	fakeData := "10.0.0.1\t10.0.0.128\t5\tGB\tMy ASN Ltd"

	// create a temp file
	fileLoc := filepath.Join(t.TempDir(), "fake-ip-data.tsv")
	err := os.WriteFile(fileLoc, []byte(fakeData), 0664)
	if err != nil {
		t.Fatalf("unable to create test file: %s", err.Error())
	}

	// run the LoadDB function
	ipData, err := LoadDb(fileLoc)
	if err != nil {
		t.Logf("function returned an error: %s", err.Error())
		t.FailNow()
	}

	testIp := ipData["10"].Child["0"].Child["0"].Child["1"].Value

	// ensure the correct data is returned
	assert.Equal(t, fakeData, fmt.Sprintf("%s\t%s\t%d\t%s\t%s",
		testIp.RangeStart,
		testIp.RangeEnd,
		testIp.ASNumber,
		testIp.CountryCode,
		testIp.ASDescription),
	)

}
