package world

import "github.com/tobz/phosphorus/utils"

type Area interface {
	AreaID() uint32
	Name() string

	Contains(utils.Point3D) bool
}

type circleArea struct {
	origin utils.Point3D
	radius float64
}

func (ca *circleArea) Name() string {
	return "Test Area"
}

func (ca *circleArea) Contains(p utils.Point3D) bool {
	return p.WithinRadius(ca.origin, ca.radius)
}
