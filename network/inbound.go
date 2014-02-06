package network

import (
	"bytes"
	"fmt"

	"github.com/tobz/phosphorus/constants"
)

type InboundPacket struct {
	basePacket
	buf    []byte
	bufPos int

	SessionID uint16
	Sequence  uint16
	Parameter uint16
}

func NewInboundPacket(source []byte, typ constants.PacketType) *InboundPacket {
	buf := make([]byte, len(source))
	copy(buf, source)

	p := &InboundPacket{basePacket{typ, constants.InvalidCode}, buf, 0, 0, 0, 0}
	p.readHeader()
	return p
}

func (p *InboundPacket) readHeader() {
	// Skip the length field because we don't care.
	p.Skip(2)

	sequence, _ := p.ReadBEUint16()
	p.Sequence = sequence

	sessionId, _ := p.ReadBEUint16()
	p.SessionID = sessionId

	parameter, _ := p.ReadBEUint16()
	p.Parameter = parameter

	code, _ := p.ReadBEUint16()
	p.code = constants.PacketCode(code)
}

func (p *InboundPacket) Skip(n int) {
	p.bufPos += n
}

func (p *InboundPacket) Buffer() []byte {
	return p.buf
}

func (p *InboundPacket) Len() int {
    return p.bufPos
}

func (p *InboundPacket) ReadUint8() (uint8, error) {
	var val uint8

	if p.canReadFurther(1) {
		val = uint8(p.buf[p.bufPos])

		p.bufPos += 1

		return val, nil
	}

	return 0, fmt.Errorf("unable to read uint8 from buffer")
}

func (p *InboundPacket) ReadUint16() (uint16, error) {
	var val uint16

	if p.canReadFurther(2) {
		val += uint16(p.buf[p.bufPos])
		val += uint16(p.buf[p.bufPos+1]) << 8

		p.bufPos += 2

		return val, nil
	}

	return 0, fmt.Errorf("unable to read uint16 from buffer")
}

func (p *InboundPacket) ReadUint32() (uint32, error) {
	var val uint32

	if p.canReadFurther(4) {
		val += uint32(p.buf[p.bufPos])
		val += uint32(p.buf[p.bufPos+1]) << 8
		val += uint32(p.buf[p.bufPos+2]) << 16
		val += uint32(p.buf[p.bufPos+3]) << 24

		p.bufPos += 4

		return val, nil
	}

	return 0, fmt.Errorf("unable to read uint32 from buffer")
}

func (p *InboundPacket) ReadUint64() (uint64, error) {
	var val uint64

	if p.canReadFurther(8) {
		val += uint64(p.buf[p.bufPos])
		val += uint64(p.buf[p.bufPos+1]) << 8
		val += uint64(p.buf[p.bufPos+2]) << 16
		val += uint64(p.buf[p.bufPos+3]) << 24
		val += uint64(p.buf[p.bufPos+4]) << 32
		val += uint64(p.buf[p.bufPos+5]) << 40
		val += uint64(p.buf[p.bufPos+6]) << 48
		val += uint64(p.buf[p.bufPos+7]) << 56

		p.bufPos += 8

		return val, nil
	}

	return 0, fmt.Errorf("unable to read uint64 from buffer")
}

func (p *InboundPacket) ReadInt8() (int8, error) {
	var val int8

	if p.canReadFurther(1) {
		val = int8(p.buf[p.bufPos])

		p.bufPos += 1

		return val, nil
	}

	return 0, fmt.Errorf("unable to read int8 from buffer")
}

func (p *InboundPacket) ReadInt16() (int16, error) {
	var val int16

	if p.canReadFurther(2) {
		val += int16(p.buf[p.bufPos])
		val += int16(p.buf[p.bufPos+1]) << 8

		p.bufPos += 2

		return val, nil
	}

	return 0, fmt.Errorf("unable to read int16 from buffer")
}

func (p *InboundPacket) ReadInt32() (int32, error) {
	var val int32

	if p.canReadFurther(4) {
		val += int32(p.buf[p.bufPos])
		val += int32(p.buf[p.bufPos+1]) << 8
		val += int32(p.buf[p.bufPos+2]) << 16
		val += int32(p.buf[p.bufPos+3]) << 24

		p.bufPos += 4

		return val, nil
	}

	return 0, fmt.Errorf("unable to read int32 from buffer")
}

func (p *InboundPacket) ReadInt64() (int64, error) {
	var val int64

	if p.canReadFurther(8) {
		val += int64(p.buf[p.bufPos])
		val += int64(p.buf[p.bufPos+1]) << 8
		val += int64(p.buf[p.bufPos+2]) << 16
		val += int64(p.buf[p.bufPos+3]) << 24
		val += int64(p.buf[p.bufPos+4]) << 32
		val += int64(p.buf[p.bufPos+5]) << 40
		val += int64(p.buf[p.bufPos+6]) << 48
		val += int64(p.buf[p.bufPos+7]) << 56

		p.bufPos += 8

		return val, nil
	}

	return 0, fmt.Errorf("unable to read int64 from buffer")
}

func (p *InboundPacket) ReadBEUint16() (uint16, error) {
	var val uint16

	if p.canReadFurther(2) {
		val += uint16(p.buf[p.bufPos+1])
		val += uint16(p.buf[p.bufPos]) << 8

		p.bufPos += 2

		return val, nil
	}

	return 0, fmt.Errorf("unable to read network-order uint16 from buffer")
}

func (p *InboundPacket) ReadBEUint32() (uint32, error) {
	var val uint32

	if p.canReadFurther(4) {
		val += uint32(p.buf[p.bufPos+3])
		val += uint32(p.buf[p.bufPos+2]) << 8
		val += uint32(p.buf[p.bufPos+1]) << 16
		val += uint32(p.buf[p.bufPos]) << 24

		p.bufPos += 4

		return val, nil
	}

	return 0, fmt.Errorf("unable to read network-order uint32 from buffer")
}

func (p *InboundPacket) ReadBEUint64() (uint64, error) {
	var val uint64

	if p.canReadFurther(8) {
		val += uint64(p.buf[p.bufPos+7])
		val += uint64(p.buf[p.bufPos+6]) << 8
		val += uint64(p.buf[p.bufPos+5]) << 16
		val += uint64(p.buf[p.bufPos+4]) << 24
		val += uint64(p.buf[p.bufPos+3]) << 32
		val += uint64(p.buf[p.bufPos+2]) << 40
		val += uint64(p.buf[p.bufPos+1]) << 48
		val += uint64(p.buf[p.bufPos]) << 56

		p.bufPos += 8

		return val, nil
	}

	return 0, fmt.Errorf("unable to read network-order uint64 from buffer")
}

func (p *InboundPacket) ReadBEInt16() (int16, error) {
	var val int16

	if p.canReadFurther(2) {
		val += int16(p.buf[p.bufPos+1])
		val += int16(p.buf[p.bufPos]) << 8

		p.bufPos += 2

		return val, nil
	}

	return 0, fmt.Errorf("unable to read network-order int16 from buffer")
}

func (p *InboundPacket) ReadBEInt32() (int32, error) {
	var val int32

	if p.canReadFurther(4) {
		val += int32(p.buf[p.bufPos+3])
		val += int32(p.buf[p.bufPos+2]) << 8
		val += int32(p.buf[p.bufPos+1]) << 16
		val += int32(p.buf[p.bufPos]) << 24

		p.bufPos += 4

		return val, nil
	}

	return 0, fmt.Errorf("unable to read network-order int32 from buffer")
}

func (p *InboundPacket) ReadBEInt64() (int64, error) {
	var val int64

	if p.canReadFurther(8) {
		val += int64(p.buf[p.bufPos+7])
		val += int64(p.buf[p.bufPos+6]) << 8
		val += int64(p.buf[p.bufPos+5]) << 16
		val += int64(p.buf[p.bufPos+4]) << 24
		val += int64(p.buf[p.bufPos+3]) << 32
		val += int64(p.buf[p.bufPos+2]) << 40
		val += int64(p.buf[p.bufPos+1]) << 48
		val += int64(p.buf[p.bufPos]) << 56

		p.bufPos += 8

		return val, nil
	}

	return 0, fmt.Errorf("unable to read network-order int64 from buffer")
}

func (p *InboundPacket) ReadBoundedString(length int) (string, error) {
	if !p.canReadFurther(1) {
		return "", fmt.Errorf("unable to read bounded string from buffer")
	}

	// set our bounds
	start := p.bufPos
	end := p.bufPos + length
	if end >= len(p.buf) {
		end = len(p.buf)
	}

	possibleStr := p.buf[start:end]
	termByte := bytes.IndexByte(possibleStr, 0x00)
	if termByte == -1 {
		// the whole thing is a string, apparently
		return string(possibleStr), nil
	}

	return string(possibleStr[:termByte]), nil
}

func (p *InboundPacket) ReadLengthPrefixedString() (string, error) {
	if !p.canReadFurther(1) {
		return "", fmt.Errorf("unable to read length-prefixed string from buffer: can't read length")
	}

	n, err := p.ReadUint8()
	if err != nil {
		return "", err
	}

	length := int(n)

	if !p.canReadFurther(length) {
		return "", fmt.Errorf("unable to read length-prefixed string from buffer: length > remaining")
	}

	start := p.bufPos
	end := p.bufPos + length

	p.bufPos += length

	return string(p.buf[start:end]), nil
}

func (p *InboundPacket) canReadFurther(n int) bool {
	return (p.bufPos + n) <= len(p.buf)
}

// No-op function to satisfy the packet interface.
func (p *InboundPacket) Finalize() {
}
