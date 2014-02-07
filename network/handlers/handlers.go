package handlers

import (
	"fmt"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/log"
	"github.com/tobz/phosphorus/network"
)

type PacketHandler func(c interfaces.Client, p *network.InboundPacket) error

type handlerMap map[constants.PacketCode]PacketHandler
type handlersMap map[constants.PacketType]handlerMap

var handlers = handlersMap{
	constants.PacketUDP: make(handlerMap),
	constants.PacketTCP: make(handlerMap),
}

func Register(typ constants.PacketType, code constants.PacketCode, handler PacketHandler) {
	if _, ok := handlers[typ][code]; ok {
		panic(fmt.Sprintf("packet handler for %v:%v is already defined!", typ, code))
	}
	handlers[typ][code] = handler
}

func Handle(c interfaces.Client, p *network.InboundPacket) error {
	packetType := "TCP"
	if p.Type() == constants.PacketUDP {
		packetType = "UDP"
	}

	// Make sure we actually have a packet handler to fulfill this request.
	if _, ok := handlers[p.Type()][p.Code()]; !ok {
		return fmt.Errorf("tried to handle packet %s(0x%X) but no registered handler found", packetType, byte(p.Code()))
	}

	log.Server.ClientDebug(c, "network", "Handling packet %s(0x%X) -> %d bytes", packetType, byte(p.Code()), p.Len())

	return handlers[p.Type()][p.Code()](c, p)
}
