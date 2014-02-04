package network

import "bytes"
import "encoding/binary"

import "github.com/tobz/phosphorus/constants"

type OutboundPacket struct {
	basePacket
	buf       *bytes.Buffer
	finalized bool
}

func NewOutboundPacket(typ constants.PacketType, code constants.PacketCode) *OutboundPacket {
	p := &OutboundPacket{basePacket{typ, code}, new(bytes.Buffer), false}

	p.buf.Write([]byte{0x00, 0x00, byte(code)})
	return p
}

func (p *OutboundPacket) WriteUint8(val uint8) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteUint16(val uint16) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteUint32(val uint32) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteUint64(val uint64) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteInt8(val int8) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteInt16(val int16) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteInt32(val int32) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteInt64(val int64) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteBEUint16(val uint16) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteBEUint32(val uint32) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteBEUint64(val uint64) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteBEInt16(val int16) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteBEInt32(val int32) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteBEInt64(val int64) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteBoundedString(val string, length int) {
	byteVal := []byte(val)
	if len(byteVal) > length {
		byteVal = byteVal[0:length]
	}

	for len(byteVal) < length {
		byteVal = append(byteVal, 0x00)
	}

	p.buf.Write(byteVal)
}

func (p *OutboundPacket) WriteRepeated(val uint8, count int) {
	buf := make([]byte, count)

	for i, _ := range buf {
		buf[i] = val
	}

	p.buf.Write(buf)
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
