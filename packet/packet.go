package packet

import "server"

type Handler func(client *server.Client, packet Packet)

type Packet interface {
	Type() Type
	Code() Code
	Buffer() []byte
}

type base struct {
	typ  Code
	code Type
}

func (p *base) Type() Type {
	return p.typ
}

func (p *base) Code() Code {
	return p.code
}
