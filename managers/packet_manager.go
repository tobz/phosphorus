package managers

import "github.com/tobz/phosphorus/constants"
import "github.com/tobz/phosphorus/interfaces"
import "github.com/tobz/phosphorus/network"

var DefaultPacketManager *PacketManager = &PacketManager{}

type PacketManager struct {
}

func (pm *PacketManager) RegisterInboundHandler(packetType constants.PacketType, packetCode constants.PacketCode, handler network.InboundPacketHandler) {
}

func (pm *PacketManager) HandlePacket(client interfaces.Client, packet *network.InboundPacket) error {
	return nil
}
