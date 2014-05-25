package world

import "os"
import "fmt"
import "sync"
import "time"
import "encoding/xml"
import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/interfaces"

type RegionsContainer struct {
	RegionEntries []RegionEntry `xml:"region"`
}

type RegionEntry struct {
	RegionID       uint16 `xml:"regionID"`
	Description    string `xml:"description"`
	DivingEnabled  bool   `xml:"isDivingEnabled"`
	WaterLevel     int32  `xml:"waterLevel"`
	Expansion      uint8  `xml:"expansion"`
	HousingEnabled bool   `xml:"isHousingEnabled"`
	Instance       bool   `xml:"instance"`
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
	var regionsContainer RegionsContainer
	err = decoder.Decode(&regionsContainer)
	if err != nil {
		return nil, err
	}

	return regionsContainer.RegionEntries, nil
}

type Region struct {
	lastMovementUpdateTime time.Time
	lastBehaviorUpdateTime time.Time

	internalZones   map[uint16]ZoneEntry
	internalAreas   map[uint32]Area
	internalObjects map[uint16]interfaces.WorldObject

	objectLock *sync.RWMutex

	tree *Octree
}

func NewRegion(re RegionEntry, zones []ZoneEntry) *Region {
	r := &Region{}

	internalAreas := make(map[uint32]Area)
	r.internalAreas = internalAreas

	internalObjects := make(map[uint16]interfaces.WorldObject)
	r.internalObjects = internalObjects

	objectLock := &sync.RWMutex{}
	r.objectLock = objectLock

	// Figure out which zones belong to this region.  Also figure out what the maximum
	// dimensions of our region need to be.
	var maxX int64
	var maxY int64

	internalZones := make(map[uint16]ZoneEntry)
	for _, zone := range zones {
		if re.RegionID == zone.RegionID {
			internalZones[zone.ZoneID] = zone

			maxZoneY := (zone.OffsetY + zone.Height)
			maxZoneX := (zone.OffsetX + zone.Width)
			if int64(maxZoneY) > maxY {
				maxY = int64(maxZoneY)
			}

			if int64(maxZoneX) > maxX {
				maxX = int64(maxZoneX)
			}
		}
	}
	r.internalZones = internalZones

	// Now create our tree.
	tree := NewOctree(maxY, maxX, 32768)
	r.tree = tree

	return r
}

func (r *Region) Tick() {
	currentTime := time.Now()

	// See if it's time to update movement.
	if currentTime.After(r.lastMovementUpdateTime.Add(constants.RegionMovementUpdateInterval)) {
		r.lastMovementUpdateTime = time.Now()
		go r.updateMovement()
	}

	// See if it's time to update behavior.
	if currentTime.After(r.lastMovementUpdateTime.Add(constants.RegionMovementUpdateInterval)) {
		r.lastBehaviorUpdateTime = time.Now()
		go r.updateBehavior()
	}
}

func (r *Region) updateMovement() {
}

func (r *Region) updateBehavior() {
}

func (r *Region) AddObject(obj interfaces.WorldObject) error {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	// Make sure they don't already exist in the region.
	if _, ok := r.internalObjects[obj.ObjectID()]; ok {
		return fmt.Errorf("object is already in region")
	}

	// Add them to the tree.
	return r.tree.AddObject(obj)
}

func (r *Region) RemoveObject(obj interfaces.WorldObject) error {
	r.objectLock.Lock()
	defer r.objectLock.Unlock()

	// Make sure they don't already exist in the region.
	if _, ok := r.internalObjects[obj.ObjectID()]; !ok {
		return fmt.Errorf("object isn't in region")
	}

	// Remove them to the tree.
	return r.tree.RemoveObject(obj)
}
