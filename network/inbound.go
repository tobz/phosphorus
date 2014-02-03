package network

import "fmt"

type InboundPacket struct {
    basePacket
	buf    []byte
	bufPos int
}

func NewInboundPacket(source []byte, typ PacketType) *InboundPacket {
	p := &Inbound{typ, InvalidCode, buf, 0}
	p.readHeader()
	return p
}

func (p *InboundPacket) readHeader() {
	// Skip the length field because we don't care.
	p.Skip(2)
	code, _ := ip.ReadUint8()
	ip.packetCode = PacketCode(code)
}

func (p *InboundPacket) Skip(n int) {
	p.bufPos += n
}

func (p *InboundPacket) Buffer() []byte {
	return p.buf
}

func (p *InboundPacket) ReadUint8() (byte, error) {
	ok := p.bufPos < len(p.buf)
	p.bufpos += 1
	if ok {
		return p.buf[p.bufPos-1], nil
	}
	return 0, fmt.Errorf("unable to read uint8 from buffer")
}
