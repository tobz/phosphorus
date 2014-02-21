package world

import "fmt"
import "sync"
import "math"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/utils"
import "github.com/tobz/phosphorus/log"

type WorldMgr struct {
	sessionQueue  interfaces.Queue
	clientMap     map[uint16]interfaces.Client
	clientMapLock *sync.RWMutex
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

	// Now load in our regions.
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
	_, err := ReadZones(zoneConfig)
	if err != nil {
		return err
	}

	// Now load our regions.
	_, err = ReadRegions(regionConfig)
	if err != nil {
		return err
	}

	return nil
}

func (w *WorldMgr) Start() error {
	log.Server.Info("world", "World manager is now running.")

	return nil
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
