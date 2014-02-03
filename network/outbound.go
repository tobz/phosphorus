package network

import "bytes"

type OutboundPacket struct {
	basePacket
	buf       *bytes.Buffer
	finalized bool
}

func NewOutboundPacket(typ PacketType, code PacketCode) *OutboundPacket {
	p := &OutboundPacket{
		typ,
		code,
		bytes.Buffer,
	}

	p.buf.Write([]byte{0x00, 0x00, byte(code)})
	return p
}

func (p *OutboundPacket) Buffer() []byte {
	if !p.finalized {
		panic("tried to get unfinalized packet content!")
	}
	return p.buf.Bytes()
}

func (p *OutboundPacket) Finalize() {
	lensize := 2
	buflen := p.buf.Len() - lensize

	b := p.buf.Bytes()
	b[0] = byte(buflen)
	b[1] = byte(buflen >> 8)

	p.finalized = true
}
