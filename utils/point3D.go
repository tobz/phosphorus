package utils

type Point3D struct {
	X float64
	Y float64
	Z float64
}

func (p *Point3D) WithinRadius(pp Point3D, r float64) bool {
    xDiff := (pp.X - p.X)
    yDiff := (pp.Y - p.Y)

    return ((xDiff * xDiff) + (yDiff * yDiff)) < (r * r)
}
