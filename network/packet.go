package network

import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/interfaces"

type PacketHandler func(client interfaces.Client, packet interfaces.Packet)

type basePacket struct {
	typ  constants.PacketType
	code constants.PacketCode
}

func (p *basePacket) Type() constants.PacketType {
	return p.typ
}

func (p *basePacket) Code() constants.PacketCode {
	return p.code
}
