package phosphorus

import "fmt"
import "bytes"

type PacketType int
type PacketCode byte

const (
    PacketType_UDP PacketType = iota
    PacketType_TCP
)

type BasePacket interface {
    Buffer() []byte
}

type InboundPacket struct {
    Type PacketType
    Code PacketCode
    buffer []byte
    bufPos int
}

type OutboundPacket struct {
    Type PacketType
    Code PacketCode
    buffer *bytes.Buffer
}

func NewInboundPacket(sourceBuf []byte, sourceLen uint64, packetType PacketType) *InboundPacket {
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

    ip.Code = ip.ReadUInt8()
}

func (ip *InboundPacket) hasNumBytes(n int) bool {
    remaining := len(ip.Buffer) - ip.bufPos
    if remaining < n {
        return false, fmt.Errorf("needed %d bytes, only have %d bytes available", n, remaining)
    }

    return true, nil
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

func NewOutboundPacket(packetType PacketType, packetCode PacketCode) *OutboundPacket {
    packet := &OutboundPacket{
        Type: packetType,
        Code: packetCode,
        buffer: &bytes.Buffer{},
    }

    // Write our packet length placeholder and the packet code.
    packet.buffer.Write([]byte{ 0x00, 0x00, packetCode })

    return packet
}

func (op *OutboundPacket) Finalize() {
    // Get the length of the buffer minus the packet length field.
    bufLength := len(op.buffer.Bytes()) - 2

    // Write in the length as a uint16 at the beginning of the buffer.
    op.buffer.Bytes()[0] = byte(bufLength)
    op.buffer.Bytes()[1] = byte(bufLength >> 8)
}
