package constants

type PacketType int
type PacketCode byte

const (
	PacketType_UDP PacketType = iota
	PacketType_TCP
)
