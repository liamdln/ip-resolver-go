package ipresolver

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveIp(t *testing.T) {

	// create a fake database file and load it
	fakeData := "10.0.0.1\t10.0.0.128\t5\tGB\tMy ASN Ltd\n10.0.0.129\t10.0.0.255\t6\tUS\tYeeHaw ASN Inc"

	// create a temp file
	fileLoc := filepath.Join(t.TempDir(), "fake-ip-data.tsv")
	err := os.WriteFile(fileLoc, []byte(fakeData), 0664)
	if err != nil {
		t.Fatalf("unable to create test file: %s", err.Error())
	}

	// load the file
	err = LoadIPFile(fileLoc)
	if err != nil {
		t.Logf("could not load IP file: %s", err.Error())
		t.FailNow()
	}

	// check the ip resolves
	data, err := ResolveIp("10.0.0.200")
	if err != nil {
		t.Logf("unable to resolve ip: %s", err.Error())
	}

	// ensure the correct asn is selected
	assert.Equal(t, "10.0.0.129", data.RangeStart)
	assert.Equal(t, "10.0.0.255", data.RangeEnd)
	assert.Equal(t, 6, data.ASNumber)
	assert.Equal(t, "US", data.CountryCode)
	assert.Equal(t, "YeeHaw ASN Inc", data.ASDescription)

}
