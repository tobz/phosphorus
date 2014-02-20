package duplicatenamecheck

import (
	"fmt"

	"github.com/tobz/phosphorus/constants"
	"github.com/tobz/phosphorus/helpers"
	"github.com/tobz/phosphorus/interfaces"
	"github.com/tobz/phosphorus/network"
	"github.com/tobz/phosphorus/network/handlers"
)

func init() {
	handlers.Register(constants.PacketTCP, constants.RequestDuplicateNameCheck, HandleDuplicateNameCheckRequest)
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
