package utils

import "testing"
import "github.com/stretchr/testify/assert"

func TestUUIDDifferentValues(t *testing.T) {
    uuidOne := NewUUID("totally random string")
    uuidTwo := NewUUID("tubular dudee")

    assert.NotEqual(t, uuidOne, uuidTwo, "UUIDs have different plaintext; should not be equal")
}

func TestUUIDSameValues(t *testing.T) {
    uuidOne := NewUUID("tubular dudee")
    uuidTwo := NewUUID("tubular dudee")

    assert.Equal(t, uuidOne, uuidTwo, "UUIDs have identical plaintext; should be equal")
}
