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

func (p *OutboundPacket) WriteUInt8(val uint8) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteHUInt16(val uint16) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteHUInt32(val uint32) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteHUInt64(val uint64) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteHInt16(val int16) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteHInt32(val int32) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteHInt64(val int64) {
	binary.Write(p.buf, binary.LittleEndian, val)
}

func (p *OutboundPacket) WriteUInt16(val uint16) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteUInt32(val uint32) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteUInt64(val uint64) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteInt16(val int16) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteInt32(val int32) {
	binary.Write(p.buf, binary.BigEndian, val)
}

func (p *OutboundPacket) WriteInt64(val int64) {
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

func (p *OutboundPacket) WriteLengthPrefixedString(val string) {
	strBytes := []byte(val)

	// We can't go over 255 characters in length.
	strLen := len(strBytes)
	if strLen > 255 {
		strLen = 255
	}

	p.WriteUInt8(uint8(strLen))
	p.buf.Write(strBytes[0:strLen])
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

func (p *OutboundPacket) Len() int {
	return p.buf.Len()
}

func (p *OutboundPacket) Finalize() {
	if p.finalized {
		return
	}

	lensize := 3
	buflen := p.buf.Len() - lensize

	b := p.buf.Bytes()
	b[0] = byte(buflen >> 8)
	b[1] = byte(buflen)

	p.finalized = true
}
