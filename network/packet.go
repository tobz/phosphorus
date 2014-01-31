package network

import "fmt"
import "bytes"
import "github.com/tobz/phosphorus/constants"

type BasePacket interface {
    Buffer() []byte
}

type InboundPacket struct {
    Type constants.PacketType
    Code constants.PacketCode
    buffer []byte
    bufPos int
}

type OutboundPacket struct {
    Type constants.PacketType
    Code constants.PacketCode
    buffer *bytes.Buffer
}

func NewInboundPacket(sourceBuf []byte, sourceLen uint64, packetType constants.PacketType) *InboundPacket {
    newBuf := make([]byte, sourceLen)
    copy(sourceBuf[:sourceLen], newBuf)

    packet := &InboundPacket{
        Type: packetType,
        buffer: newBuf,
    }

    packet.readHeader()

    return packet
}

func (ip *InboundPacket) readHeader() {
    // Skip the length field because we don't care.
    ip.Skip(2)

    code, _ := ip.ReadUInt8()
    ip.Code = PacketCode(code)
}

func (ip *InboundPacket) hasNumBytes(n int) (bool, error) {
    remaining := len(ip.buffer) - ip.bufPos
    if remaining < n {
        return false, fmt.Errorf("needed %d bytes, only have %d bytes available", n, remaining)
    }

    return true, nil
}

func (op *InboundPacket) Buffer() []byte {
    return op.buffer
}

func (ip *InboundPacket) Skip(n int) {
    ip.bufPos += n
}

func (ip *InboundPacket) ReadUInt8() (byte, error) {
    if ok, err := ip.hasNumBytes(1); !ok {
        return 0, err
    }

    ip.bufPos++
    return ip.buffer[ip.bufPos - 1], nil
}

func NewOutboundPacket(packetType constants.PacketType, packetCode constants.PacketCode) *OutboundPacket {
    packet := &OutboundPacket{
        Type: packetType,
        Code: packetCode,
        buffer: &bytes.Buffer{},
    }

    // Write our packet length placeholder and the packet code.
    packet.buffer.Write([]byte{ 0x00, 0x00, byte(packetCode) })

    return packet
}

func (op *OutboundPacket) Buffer() []byte {
    return op.buffer.Bytes()
}

func (op *OutboundPacket) Finalize() {
    // Get the length of the buffer minus the packet length field.
    bufLength := len(op.buffer.Bytes()) - 2

    // Write in the length as a uint16 at the beginning of the buffer.
    op.buffer.Bytes()[0] = byte(bufLength)
    op.buffer.Bytes()[1] = byte(bufLength >> 8)
}
