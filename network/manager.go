package phosphorus

var DefaultPacketManager *PacketManager = &PacketManager{}

type PacketManager struct {
}

func (pm *PacketManager) RegisterInboundHandler(handler *InboundPacketHandler) {
}

func (pm *PacketManager) Send(client *Client, handler *OutboundPacketHandler) {
}
