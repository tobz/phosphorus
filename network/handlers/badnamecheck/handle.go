package badnamecheck

import (
	"fmt"
	"strings"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestBadNameCheck, HandleBadNameCheckRequest)
}

func HandleBadNameCheckRequest(c interfaces.Client, p *network.InboundPacket) error {
	name, err := p.ReadBoundedString(30)
	valid := true

	invalidWords, err := c.Server().Config().GetAsManyStrings("server/invalidWords")
	if err != nil {
		return err
	}

	for _, invalidWord := range invalidWords {
		if strings.Contains(name, invalidWord) {
			valid = false
			break
		}
	}

	return SendBadNameCheckResponse(c, name, valid)
}

func SendBadNameCheckResponse(c interfaces.Client, name string, isValid bool) error {
	p := network.NewOutboundPacket(constants.PacketTCP, constants.ResponseBadNameCheck)
	p.WriteBoundedString(name, 30)

	if c.Account() != nil {
		return fmt.Errorf("client didn't have an account!")
	}

	p.WriteRepeated(0x00, 20)
	if isValid {
		p.WriteUint8(0x01)
	} else {
		p.WriteUint8(0x00)
	}

	return c.Send(p)
}
