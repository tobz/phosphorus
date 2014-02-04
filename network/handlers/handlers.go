package handlers

import (
	"fmt"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
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
	return handlers[p.Type()][p.Code()](c, p)
}
