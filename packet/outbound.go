package packet

import "bytes"

type Outbound struct {
	base
	buf       *bytes.Buffer
	finalized bool
}

func NewOutbound(typ Type, code Code) *Outbound {
	p := &Outbound{
		typ,
		code,
		bytes.Buffer,
	}

	p.buf.Write([]byte{0x00, 0x00, byte(code)})
	return p
}

func (p *Outbound) Buffer() []byte {
	if !p.finalized {
		panic("Tried to get unfinalized packet content!")
	}
	return p.buf.Bytes()
}

func (p *Outbound) Finalize() {
	lensize := 2
	buflen := p.buf.Len() - lensize

	b := p.buf.Bytes()
	b[0] = byte(buflen)
	b[1] = byte(buflen >> 8)

	p.finalized = true
}
