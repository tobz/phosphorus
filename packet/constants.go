package packet

type Type int

const (
	UDP Type = iota
	TCP
)

type Code byte

const (
	InvalidCode Code = 0x0

	RequestNameCheck  Code = 0x6A ^ 168
	ResponseNameCheck Code = 0xC3
)
