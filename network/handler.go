package network

import "github.com/tobz/phosphorus/interfaces"

type InboundPacketHandler func(client interfaces.Client, packet *InboundPacket) error
