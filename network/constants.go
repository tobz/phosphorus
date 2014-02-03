package network

type PacketType int

const (
	PacketUDP PacketType = iota
	PacketTCP
)

type PacketCode byte

const (
	InvalidCode PacketCode = 0x0
)
