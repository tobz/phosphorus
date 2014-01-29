package network

import "github.com/tobz/phosphorus"

var DefaultPacketManager *PacketManager = &PacketManager{}

type PacketManager struct {
}

func (pm *PacketManager) RegisterInboundHandler(handler *InboundPacketHandler) {
}

func (pm *PacketManager) Send(client *phosphorus.Client, handler *OutboundPacketHandler) {
}
