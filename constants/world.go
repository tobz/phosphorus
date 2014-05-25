package constants

import "time"

const (
	DefaultOctreeDepth             = 6
	WorldManagerRegionTickInterval = time.Millisecond * 200
	WorldManagerUpdateTickInterval = time.Second
	RegionMovementUpdateInterval   = time.Millisecond * 500
	RegionBehaviorUpdateInterval   = time.Millisecond * 250
)
