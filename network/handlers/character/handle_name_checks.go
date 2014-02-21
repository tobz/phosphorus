package character

import (
	"fmt"
	"strings"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/helpers"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestBadNameCheck, HandleBadNameCheckRequest)
	handlers.Register(constants.PacketTCP, constants.RequestDuplicateNameCheck, HandleDuplicateNameCheckRequest)
}

func HandleBadNameCheckRequest(c interfaces.Client, p *network.InboundPacket) error {
	name, err := p.ReadBoundedString(30)
	if err != nil {
		return err
	}

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
		p.WriteUInt8(0x01)
	} else {
		p.WriteUInt8(0x00)
	}

	return c.Send(p)
}

func HandleDuplicateNameCheckRequest(c interfaces.Client, p *network.InboundPacket) error {
	characterName, err := p.ReadBoundedString(30)
	if err != nil {
		return err
	}

	nameTaken, err := helpers.IsCharacterNameTaken(c, characterName)
	if err != nil {
		return err
	}

	return SendDuplicateNameCheckResponse(c, characterName, nameTaken)
}

func SendDuplicateNameCheckResponse(c interfaces.Client, characterName string, nameExists bool) error {
	if c.Account() == nil {
		return fmt.Errorf("client has no account!")
	}

	p := network.NewOutboundPacket(constants.PacketTCP, constants.ResponseDuplicateNameCheck)

	p.WriteBoundedString(characterName, 30)
	p.WriteBoundedString(c.Account().Username, 24)

	if nameExists {
		p.WriteUInt8(0x01)
	} else {
		p.WriteUInt8(0x00)
	}

	p.WriteRepeated(0x00, 3)

	return c.Send(p)
}
