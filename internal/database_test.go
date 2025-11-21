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
	fakeData := "10.0.0.1\t10.0.0.255\t5\tGB\tMy ASN Inc"

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

	// ensure the correct data is returned
	assert.Equal(t, fakeData, fmt.Sprintf("%s\t%s\t%d\t%s\t%s",
		ipData[0].RangeStart,
		ipData[0].RangeEnd,
		ipData[0].ASNumber,
		ipData[0].CountryCode,
		ipData[0].ASDescription),
	)

}
