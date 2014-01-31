package network

import "github.com/tobz/phosphorus/interfaces"

type BasePacketHandler struct {
    packetType PacketType
    packetCode uint64
}

func (bh *BasePacketHandler) GetPacketType() PacketType {
    return bh.packetType
}

func (bh *BasePacketHandler) GetPacketCode() uint64 {
    return bh.packetCode
}

type BasicPacketHandler interface {
    GetPacketType() PacketType
    GetPacketCode() uint64
}

type InboundPacketHandler func(client *interfaces.Client, packet *InboundPacket) error
