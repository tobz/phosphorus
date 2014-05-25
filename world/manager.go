package world

import "fmt"
import "sync"
import "math"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/utils"
import "github.com/tobz/phosphorus/log"
import "github.com/tobz/phosphorus/timers"
import "github.com/tobz/phosphorus/constants"

type WorldMgr struct {
	sessionQueue  interfaces.Queue
	clientMap     map[uint16]interfaces.Client
	clientMapLock *sync.RWMutex

	regions           map[uint16]*Region
	regionUpdateTimer *timers.Timer

	worldUpdateTimer *timers.Timer
}

func NewWorldMgr(c interfaces.Config) (*WorldMgr, error) {
	// Get the player limit and make sure it's in the right range.
	playerLimit, err := c.GetAsInteger("server/playerLimit")
	if err != nil {
		return nil, err
	}

	if playerLimit < int64(1) || playerLimit > int64(math.MaxUint16) {
		return nil, fmt.Errorf("server/playerLimit must be from 0 to 65535!")
	}

	w := &WorldMgr{}
	w.sessionQueue = utils.NewQueue()
	w.clientMap = make(map[uint16]interfaces.Client)
	w.clientMapLock = &sync.RWMutex{}

	// Seed our session ID queue.
	for i := int64(0); i < playerLimit; i++ {
		w.sessionQueue.Push(uint16(i))
	}

	// Set up our timers.
	w.regionUpdateTimer = timers.NewTrackedTimer("regionUpdateTimer", constants.WorldManagerRegionTickInterval)

	w.worldUpdateTimer = timers.NewTrackedTimer("worldUpdateTime", constants.WorldManagerUpdateTickInterval)
	w.worldUpdateTimer.AddSink(w)

	// Now load in our regions.
	w.regions = make(map[uint16]*Region)

	regionConfig, err := c.GetAsString("world/regions")
	if err != nil {
		return nil, fmt.Errorf("couldn't find the location of the region data file!")
	}

	zoneConfig, err := c.GetAsString("world/zones")
	if err != nil {
		return nil, fmt.Errorf("couldn't find the location of the zone data file!")
	}

	err = w.loadRegionsAndZones(regionConfig, zoneConfig)
	if err != nil {
		return nil, err
	}

	return w, nil
}

func (w *WorldMgr) loadRegionsAndZones(regionConfig, zoneConfig string) error {
	// Load up our zones first.  We need them to properly size our regions.
	zoneEntries, err := ReadZones(zoneConfig)
	if err != nil {
		return err
	}

	// Now load our regions.
	regionEntries, err := ReadRegions(regionConfig)
	if err != nil {
		return err
	}

	// Now go through each region entry and create it.
	for _, regionEntry := range regionEntries {
		region := NewRegion(regionEntry, zoneEntries)

		// Register the region for tick updates.
		w.regionUpdateTimer.AddSink(region)

		// Store the region.
		if _, ok := w.regions[regionEntry.RegionID]; ok {
			return fmt.Errorf("tried to register region #%d, but region ID already register!", regionEntry.RegionID)
		}

		w.regions[regionEntry.RegionID] = region
	}

	return nil
}

func (w *WorldMgr) Start() error {
	// Start our update timers which will start driving all the behavior in the game world.
	w.regionUpdateTimer.Start()
	w.worldUpdateTimer.Start()

	log.Server.Info("world", "World manager is now running.")

	return nil
}

func (w *WorldMgr) Tick() {
	// Looking for players who haven't sent their ping heartbeat in a while and disconnect them.

	// Perform XYZ task, ABC task, etc.
}

func (w *WorldMgr) AddClient(c interfaces.Client) error {
	w.clientMapLock.Lock()
	defer w.clientMapLock.Unlock()

	v := w.sessionQueue.Pop()
	if v == nil {
		return fmt.Errorf("failed to get a session ID: no IDs left.  Server limit reached?")
	}

	sessionId, ok := v.(uint16)
	if !ok {
		return fmt.Errorf("got back a non-uint16 from the session ID pool!")
	}

	if c2, ok := w.clientMap[sessionId]; ok {
		return fmt.Errorf("grabbed new session ID #%d, but client already mapped! client: %#v", sessionId, c2)
	}

	w.clientMap[sessionId] = c
	c.SetSessionID(sessionId)

	return nil
}

func (w *WorldMgr) RemoveClient(c interfaces.Client) error {
	w.clientMapLock.Lock()
	defer w.clientMapLock.Unlock()

	if _, ok := w.clientMap[c.SessionID()]; !ok {
		return fmt.Errorf("client doesn't exist in map! client: %#v, session ID: %d", c, c.SessionID())
	}

	delete(w.clientMap, c.SessionID())
	w.sessionQueue.Push(c.SessionID())

	return nil
}
