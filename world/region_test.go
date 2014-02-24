package world

import "testing"
import "io/ioutil"
import "strings"
import "github.com/stretchr/testify/assert"

const (
	RegionFileNonExistent = "/tmp/32432e2783uige78uig3827g7g878u87ug87gh3"
	RegionFileRandom      = "/tmp/1323wed4332178egwyuu86732yugf78sd123e32"
	RegionFileValid       = "/tmp/6348912uh782h9jinh8huir89iohwef892hry98"
)

func TestLoadingNonExistentRegionFile(t *testing.T) {
	regions, err := ReadRegions(RegionFileNonExistent)
	assert.Nil(t, regions)
	assert.NotNil(t, err)
}

func TestLoadingInvalidRegionFile(t *testing.T) {
	// Some random data that will NOT parse as XML.
	randomData := []byte("#@324dfs3412e1qwde1$!#$!#@$")

	// Write our test file.
	err := ioutil.WriteFile(RegionFileRandom, randomData, 0644)
	assert.Nil(t, err)

	// Now try to read in that file as our regions file.
	regions, err := ReadRegions(RegionFileRandom)
	assert.Nil(t, regions)
	assert.NotNil(t, err)
}

func TestLoadingValidRegionFile(t *testing.T) {
	// Some random data that will NOT parse as XML.
	randomData := `
    <?xml version="1.0" encoding="utf-8"?>
    <root>
        <region>
            <regionID>789</regionID>
            <description>Test Region</description>
            <divingEnabled>True</divingEnabled>
            <waterLevel>2</waterLevel>
            <expansion>4</expansion>
            <isHousingEnabled>False</isHousingEnabled>
            <instance>False</instance>
        </region>
    </root>`
	randomData = strings.TrimSpace(randomData)

	// Write our test file.
	err := ioutil.WriteFile(RegionFileValid, []byte(randomData), 0644)
	assert.Nil(t, err)

	// Now try to read in that file as our zones file.
	regions, err := ReadRegions(RegionFileValid)
	assert.NotNil(t, regions)
	assert.Nil(t, err)

	// We should only have one region.
	assert.Equal(t, 1, len(regions), "there should be 1 region")

	// Test two of the values we know should be specific in there.
	assert.Equal(t, uint32(789), regions[0].RegionID, "region ID doesn't match")
	assert.Equal(t, uint32(4), regions[0].Expansion, "expansion doesn't match")
}
