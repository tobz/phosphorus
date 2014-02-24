package utils

import "testing"
import "github.com/stretchr/testify/assert"

func TestPacketChecksum(t *testing.T) {
    // This is an actual 0xF4 packet from a 1.114 client sans the checksum at the end.
    testData := []byte{0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xF4, 0x00, 0x36, 0x01, 0x0B, 0x04}

    checksum := CalculatePacketChecksum(testData, 0, len(testData))
    assert.Equal(t, uint16(0x70D3), checksum, "checksums should be equal")
}
