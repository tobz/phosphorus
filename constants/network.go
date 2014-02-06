package constants

type PacketType int

const (
	PacketUDP PacketType = iota
	PacketTCP
)

type PacketCode byte

const (
	InvalidCode          PacketCode = 0x0
	RequestCryptKey      PacketCode = 0xF4
	ResponseCryptKey     PacketCode = 0x22
	RequestBadNameCheck  PacketCode = 0xC2
	ResponseBadNameCheck PacketCode = 0xC3
	RequestLogin         PacketCode = 0xA7
	OneWayLoginDenied    PacketCode = 0x2C
    OneWayLoginGranted   PacketCode = 0x2A
)

const (
	InboundPacketHeaderSize = 12
)
