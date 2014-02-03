package handlers

import (
	"fmt"

	"github.com/tobz/phosphorus/network"
)

type handlerMap map[network.PacketCode]network.PacketHandler
type handlersMap map[network.PacketType]HandlerMap

var handlers = handlersMap{
    network.PacketUDP: make(handlerMap),
	network.PacketTCP: make(handlerMap),
}

func Register(typ network.PacketType, code network.PacketCode, handler network.PacketHandler) {
	if _, ok := handlers[typ][code]; ok {
		panic(fmt.Sprintf("packet handler for %v:%v is already defined!", typ, code))
	}
	handlers[typ][code] = handler
}

func Handle(c *server.Client, p network.Packet) error {
	return handlers[p.Type()][p.Code()](c, p)
}
