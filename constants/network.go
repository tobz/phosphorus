package constants

type PacketType int

const (
	PacketUDP PacketType = iota
	PacketTCP
)

type PacketCode byte

const (
	InvalidCode          PacketCode = 0x0
	RequestBadNameCheck  PacketCode = 0x6A ^ 168
	ResponseBadNameCheck PacketCode = 0xC3
)

const (
	InboundPacketHeaderSize = 12
)
