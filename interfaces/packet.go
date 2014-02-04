package interfaces

import "github.com/tobz/phosphorus/constants"

type Packet interface {
	Type() constants.PacketType
	Code() constants.PacketCode
	Buffer() []byte
}
