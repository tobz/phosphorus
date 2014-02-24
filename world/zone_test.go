package world

import "testing"
import "io/ioutil"
import "strings"
import "github.com/stretchr/testify/assert"

const (
	ZoneFileNonExistent = "/tmp/32432e2783uige78uig3827g7g878u87ug87gh3"
	ZoneFileRandom      = "/tmp/1323wed4332178egwyuu86732yugf78sd123e32"
	ZoneFileValid       = "/tmp/6348912uh782h9jinh8huir89iohwef892hry98"
)

func TestLoadingNonExistentZoneFile(t *testing.T) {
	zones, err := ReadZones(ZoneFileNonExistent)
	assert.Nil(t, zones)
	assert.NotNil(t, err)
}

func TestLoadingInvalidZoneFile(t *testing.T) {
	// Some random data that will NOT parse as XML.
	randomData := []byte("#@324dfs3412e1qwde1$!#$!#@$")

	// Write our test file.
	err := ioutil.WriteFile(ZoneFileRandom, randomData, 0644)
	assert.Nil(t, err)

	// Now try to read in that file as our zones file.
	zones, err := ReadZones(ZoneFileRandom)
	assert.Nil(t, zones)
	assert.NotNil(t, err)
}

func TestLoadingValidZoneFile(t *testing.T) {
	// Some random data that will NOT parse as XML.
	randomData := `
    <?xml version="1.0" encoding="utf-8"?>
    <root>
        <zone>
            <zoneID>456</zoneID>
            <regionID>789</regionID>
            <description>Test Zone</description>
            <height>8</height>
            <width>8</width>
            <offsetx>8</offsetx>
            <offsety>8</offsety>
        </zone>
    </root>`
	randomData = strings.TrimSpace(randomData)

	// Write our test file.
	err := ioutil.WriteFile(ZoneFileValid, []byte(randomData), 0644)
	assert.Nil(t, err)

	// Now try to read in that file as our zones file.
	zones, err := ReadZones(ZoneFileValid)
	assert.NotNil(t, zones)
	assert.Nil(t, err)

	// We should only have one zone.
	assert.Equal(t, 1, len(zones), "there should be 1 zone")

	// Test two of the values we know should be specific in there.
	assert.Equal(t, uint32(456), zones[0].ZoneID, "zone ID doesn't match")
	assert.Equal(t, uint32(789), zones[0].RegionID, "region ID doesn't match")
}
