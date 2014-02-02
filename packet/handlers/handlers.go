package handlers

import (
	"fmt"

	"github.com/tobz/phosphorus/packet"
)

type handlermap map[packet.Code]packet.Handler
type handlersmap map[packet.Type]handlermap

var handlers = handlersmap{
	packet.UDP: make(handlermap),
	packet.TCP: make(handlermap),
}

func Register(typ packet.Type, code packet.Code, handler packet.Handler) {
	if _, ok := handlers[typ][code]; ok {
		panic(fmt.Sprintf("Packet handler for %v:%v is already defined!", typ, code))
	}
	handlers[typ][code] = handler
}

func Handle(c *server.Client, p packet.Packet) error {
	return handlers[p.Type][p.Code](c, p)
}

func All() handlersmap {
	return handlers
}
