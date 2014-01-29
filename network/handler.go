package network

import "github.com/tobz/phosphorus"

type BaseHandler struct {
    packetType PacketType
    packetCode uint64
}

func (bh *BaseHandler) GetPacketType() PacketType {
    return bh.packetType
}

func (bh *BaseHandler) GetPacketCode() uint64 {
    return bh.packetCode
}

type BasePacketHandler interface {
    GetPacketType() PacketType
    GetPacketCode() uint64
}

type InboundPacketHandler interface {
    BasePacketHandler
    HandleRequest(client *phosphorus.Client, packet *phosphorus.Packet) error
}

type OutboundPacketHandler interface {
    BasePacketHandler
    SendResponse(client *phosphorus.Client) error
}
