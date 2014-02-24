package world

import "os"
import "encoding/xml"

type ZonesContainer struct {
	ZoneEntries []ZoneEntry `xml:"zone"`
}

type ZoneEntry struct {
	ZoneID      uint32 `xml:"zoneID"`
	RegionID    uint32 `xml:"regionID"`
	Description string `xml:"description"`
	Height      uint32 `xml:"height"`
	Width       uint32 `xml:"width"`
	OffsetX     uint32 `xml:"offsetx"`
	OffsetY     uint32 `xml:"offsety"`
}

func ReadZones(zoneConfig string) ([]ZoneEntry, error) {
	// Try and open our zone config file.
	f, err := os.Open(zoneConfig)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Create a new XML decoder.
	decoder := xml.NewDecoder(f)

	// Now kith.
	var zonesContainer ZonesContainer
	err = decoder.Decode(&zonesContainer)
	if err != nil {
		return nil, err
	}

	return zonesContainer.ZoneEntries, nil
}
