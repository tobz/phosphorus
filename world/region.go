package world

import "os"
import "encoding/xml"

type RegionContainer struct {
    RegionEntries []RegionEntry `xml:"Region"`
}

type RegionEntry struct {
    RegionID uint32 `xml:"id"`
    Description string `xml:"description"`
    DivingEnabled bool `xml:"isDivingEnabled"`
    WaterLevel int32 `xml:"waterLevel"`
    Expansion uint32 `xml:"expansion"`
    HousingEnabled bool `xml:"isHousingEnabled"`
    Instance bool `xml:"instance"`
}

func ReadRegions(regionConfig string) ([]RegionEntry, error) {
    // Try and open our region config file.
    f, err := os.Open(regionConfig)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    // Create a new XML decoder.
    decoder := xml.NewDecoder(f)

    // Now kith.
    var regionContainer RegionContainer
    err = decoder.Decode(&regionContainer)
    if err != nil {
        return nil, err
    }

    return regionContainer.RegionEntries, nil
}
