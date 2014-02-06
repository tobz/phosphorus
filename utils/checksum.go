package utils

func CalculatePacketChecksum(buf []byte, offset, length int) uint16 {
	one := uint8(0x7E)
	two := uint8(0x7E)

    for i := 0; i < length; i++ {
		one += buf[offset + i]
		two += one
	}

	return uint16(two) - ((uint16(one) + uint16(two)) << 8)
}
