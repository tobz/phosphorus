package utils

import "hash/crc32"

var tableCastagnoli = crc32.MakeTable(crc32.Castagnoli)

func NewUUID(s string) uint32 {
	return crc32.Checksum([]byte(s), tableCastagnoli)
}
