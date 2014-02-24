package utils

import "testing"
import "github.com/stretchr/testify/assert"

func TestPointWithinRadius(t *testing.T) {
    p := Point3D{0, 0, 0}
    pp := Point3D{100, 100, 100}

    assert.True(t, pp.WithinRadius(p, 300))
}

func TestPointNotWithinRadius(t *testing.T) {
    p := Point3D{0, 0, 0}
    pp := Point3D{100, 100, 100}

    assert.False(t, pp.WithinRadius(p, 50))
}
