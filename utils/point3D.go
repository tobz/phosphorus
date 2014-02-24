package utils

type Point3D struct {
	X int64
	Y int64
	Z int64
}

func (p *Point3D) WithinRadius(pp Point3D, r int64) bool {
	xDiff := (pp.X - p.X)
	yDiff := (pp.Y - p.Y)

	return ((xDiff * xDiff) + (yDiff * yDiff)) < (r * r)
}
