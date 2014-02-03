package namecheck

import (
	"strings"

	"github.com/tobz/phosphorus/packet"
	"github.com/tobz/phosphorus/packet/handlers"
	"github.com/tobz/phosphorus/server"
)

func init() {
	handlers.Register(packet.TCP, packet.RequestNameCheck, nameCheck)
}

func nameCheck(client *server.Client, p *packet.Inbound) error {
	name := p.ReadString(30)
	valid := true

	for _, word := range disallowedWords {
		if strings.Contains(name, word) {
			valid = false
			break
		}
	}

	return respond(client, name, valid)
}

func respond(client *server.Client, name string, isValid bool) error {
	p := packet.NewOutbound(packet.TCP, packet.RespondNameCheck)
	p.WriteBoundedString(name, 30)

	if client.Account != nil {
		panic("Client didn't have an account!")
	}

	p.WriteRepeated(0x00, 20)
	if isValid {
		p.WriteUint8(0x01)
	} else {
		p.WriteUint8(0x00)
	}

	return client.Send(p)
}
