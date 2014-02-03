package network

import "github.com/tobz/phosphorus/server"

type PacketHandler func(client *server.Client, packet Packet)

type Packet interface {
	Type() PacketType
	Code() PacketCode
	Buffer() []byte
}

type basePacket struct {
	typ  PacketType
	code PacketCode
}

func (p *basePacket) Type() PacketType {
	return p.typ
}

func (p *basePacket) Code() PacketCode {
	return p.code
}
