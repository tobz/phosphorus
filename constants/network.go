package constants

type PacketType int

const (
	PacketUDP PacketType = iota
	PacketTCP
)

type PacketCode byte

const (
	InvalidCode                PacketCode = 0x00
	RequestPing                PacketCode = 0xA3
	ResponsePing               PacketCode = 0x29
	RequestCryptKey            PacketCode = 0xF4
	ResponseCryptKey           PacketCode = 0x22
	RequestLogin               PacketCode = 0xA7
	ServerOnlyLoginDenied      PacketCode = 0x2C
	ServerOnlyLoginGranted     PacketCode = 0x2A
	RequestCharacterSelect     PacketCode = 0x10
	ServerOnlySessionID        PacketCode = 0x28
	RequestCharacterOverview   PacketCode = 0xFC
	ResponseCharacterOverview  PacketCode = 0xFD
	ServerOnlyRealm            PacketCode = 0xFE
	ClientOnlyRealm            PacketCode = 0xAC
	RequestBadNameCheck        PacketCode = 0xC2
	ResponseBadNameCheck       PacketCode = 0xC3
	RequestDuplicateNameCheck  PacketCode = 0xCB
	ResponseDuplicateNameCheck PacketCode = 0xCC
	RequestCharacterCreate     PacketCode = 0xFF
)

const (
	InboundPacketHeaderSize = 12
)
