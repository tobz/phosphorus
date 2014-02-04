package utils

func CalculatePacketChecksum(buf []byte, offset, length int) uint16 {
	one := uint8(0x7E)
	two := uint8(0x7E)

	for offset < (offset + length) {
		one += buf[offset]
		two += one

		offset += 1
	}

	return uint16(two - ((one + two) << 8))
}
