package packet

import "fmt"

type Inbound struct {
	base
	buf    []byte
	bufpos int
}

func NewInbound(source []byte, typ Type) *Inbound {
	p := &Inbound{typ, InvalidCode, buf, 0}
	p.readHeader()
	return p
}

func (p *Inbound) readHeader() {
	// Skip the length field because we don't care.
	p.Skip(2)
	code, _ := ip.ReadUint8()
	ip.packetCode = Code(code)
}

func (p *Inbound) Skip(n int) {
	p.bufpos += n
}

func (p *Inbound) Buffer() []byte {
	return p.buf
}

func (p *Inbound) ReadUint8() (byte, error) {
	ok := p.bufpos < len(p.buf)
	p.bufpos += 1
	if ok {
		return p.buf[p.bufpos-1], nil
	}
	return 0, fmt.Errorf("Unable to read uint8 from buffer.")
}
